package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users-consumer/model/args"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/code"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users-consumer/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users-consumer/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
)

const (
	templateUserRegister = "用户: +%v-%v, 于%v 注册成功, 状态为: %v"
)

func UserRegisterNoticeConsume(ctx context.Context, body string) error {
	// 通知消息解码
	var businessMsg args.CommonBusinessMsg
	var err error
	err = json.Unmarshal(body, &businessMsg)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "body:%v Unmarshal err: %v", body, err)
		return err
	}
	if businessMsg.Type != args.UserStateEventTypeRegister {
		return nil
	}
	var notice args.UserRegisterNotice
	err = json.Unmarshal(businessMsg.Msg, &notice)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "businessMsg.Msg: %v Unmarshal err: %v", businessMsg.Msg, err)
		return err
	}
	// 获取用户信息
	serverName := args.RpcServiceMicroMallUsers
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	defer conn.Close()

	client := users.NewUsersServiceClient(conn)
	r := users.GetUserInfoByPhoneRequest{
		CountryCode: notice.CountryCode,
		Phone:       notice.Phone,
	}
	rsp, err := client.GetUserInfoByPhone(ctx, &r)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfoByPhone %v,err: %v, r: %+v", serverName, err, r)
		return err
	}
	if rsp.Common.Code != users.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfoByPhone %v,not ok : %v, rsp: %+v", serverName, err, rsp)
		return fmt.Errorf(rsp.Common.Msg)
	}
	if rsp.Info == nil || rsp.Info.AccountId == "" {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfoByPhone %v,accountId null : %v, rsp: %+v", serverName, err, rsp)
		return fmt.Errorf(errcode.GetErrMsg(code.UserNotExist))
	}

	// 为用户初始化账户
	r2 := users.CreateUserAccountRequest{
		Uid:          rsp.Info.Uid,
		Balance:      "100.00",
		CoinType:     users.CoinType_CNY,
		CoinDesc:     "CNY",
		AccountState: users.AccountState_NORMAL,
		AccountType:  users.AccountType_PERSONAL,
	}
	rsp2, err := client.CreateUserAccount(ctx, &r2)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateUserAccount %v,err: %v, r: %+v", serverName, err, r2)
		return err
	}
	if rsp2.Common.Code != users.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "GetUserInfoByPhone %v,not ok : %v, rsp: %+v", serverName, err, rsp)
		if rsp2.Common.Code == users.RetCode_ERROR {
			return fmt.Errorf(rsp.Common.Msg)
		} else if rsp2.Common.Code == users.RetCode_USER_NOT_EXIST {
			return fmt.Errorf(rsp.Common.Msg)
		}
	}
	// 发送注册成功邮件
	emailNotice := fmt.Sprintf(templateUserRegister, notice.CountryCode, notice.Phone, notice.Time, notice.State)
	err = email.SendEmailNotice(ctx, "565608463@qq.com", vars.AppName, emailNotice)

	if err != nil {
		kelvins.ErrLogger.Info(ctx, "SendEmailNotice err, emailNotice: %v", emailNotice)
		return err
	}

	return nil
}

func UserRegisterNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
