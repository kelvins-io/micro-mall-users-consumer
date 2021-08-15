package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users-consumer/model/args"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users-consumer/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-users-consumer/repository"
	"gitee.com/cristiane/micro-mall-users-consumer/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"time"
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
	time.Sleep(3 * time.Second) // 注册事务先提交
	// 获取用户信息
	user, err := repository.GetUserByPhone("id,account_id", notice.CountryCode, notice.Phone)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone ,err: %v, req: %+v", err, notice)
		return err
	}
	if user.Id <= 0 {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone user not exist %v-%v user: %+v", notice.CountryCode, notice.Phone, user)
		return fmt.Errorf("user not exist")
	}
	// 发送注册成功邮件
	emailNotice := fmt.Sprintf(args.UserRegisterTemplate, notice.CountryCode, notice.Phone, notice.Time, args.UserStateText[notice.State])
	for _, receiver := range vars.EmailNoticeSetting.Receivers {
		err = email.SendEmailNotice(ctx, receiver, vars.AppName, emailNotice)
		if err != nil {
			kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
		}
	}

	// 为用户初始化账户
	serverName := args.RpcServiceMicroMallPay
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return err
	}
	defer conn.Close()
	balanceInit := "1.9999"
	accountReq := pay_business.CreateAccountRequest{
		Owner:       user.AccountId,
		AccountType: pay_business.AccountType_Person,
		CoinType:    pay_business.CoinType_CNY,
		Balance:     balanceInit,
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
		case pay_business.RetCode_ERROR, pay_business.RetCode_USER_NOT_EXIST:
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount %v,not ok : %v, rsp: %+v", serverName, err, accountRsp)
			return fmt.Errorf(accountRsp.Common.Msg)
		}
	}

	// 发送初始个人账户成功邮件
	emailNotice = fmt.Sprintf(args.UserCreateAccountTemplate, notice.CountryCode, notice.Phone, notice.Time, balanceInit)
	for _, receiver := range vars.EmailNoticeSetting.Receivers {
		err = email.SendEmailNotice(ctx, receiver, vars.AppName, emailNotice)
		if err != nil {
			kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
		}
	}

	return nil
}

func UserRegisterNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
