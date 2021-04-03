package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users-consumer/model/args"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-users-consumer/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-users-consumer/repository"
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
	user, err := repository.GetUserByPhone("account_id", notice.CountryCode, notice.Phone)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone ,err: %v, req: %+v", err, notice)
		return err
	}
	if user.AccountId == "" {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone AccountId empty")
		return fmt.Errorf("user AccountId empty")
	}
	// 为用户初始化账户
	serverName := args.RpcServiceMicroMallPay
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	defer conn.Close()
	accountReq := pay_business.CreateAccountRequest{
		Owner:       user.AccountId,
		AccountType: pay_business.AccountType_Person,
		CoinType:    pay_business.CoinType_CNY,
		Balance:     "999999999999999999.9999",
	}
	client := pay_business.NewPayBusinessServiceClient(conn)
	accountRsp, err := client.CreateAccount(ctx, &accountReq)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateAccount %v,err: %v, r: %+v", serverName, err, accountReq)
		return err
	}
	if accountRsp.Common.Code != pay_business.RetCode_SUCCESS {
		switch accountRsp.Common.Code {
		case pay_business.RetCode_ACCOUNT_EXIST:
			return nil
		case pay_business.RetCode_ERROR,pay_business.RetCode_USER_NOT_EXIST:
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount %v,not ok : %v, rsp: %+v", serverName, err, accountRsp)
			return fmt.Errorf(accountRsp.Common.Msg)
		}
	}
	// 发送注册成功邮件
	//emailNotice := fmt.Sprintf(templateUserRegister, notice.CountryCode, notice.Phone, notice.Time, notice.State)
	//err = email.SendEmailNotice(ctx, "565608463@qq.com", vars.AppName, emailNotice)
	//
	//if err != nil {
	//	kelvins.ErrLogger.Info(ctx, "SendEmailNotice err, emailNotice: %v", emailNotice)
	//	return err
	//}

	return nil
}

func UserRegisterNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
