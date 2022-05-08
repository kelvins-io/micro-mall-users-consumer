package service

import (
	"context"
	"fmt"

	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util"
	"gitee.com/cristiane/micro-mall-users-consumer/proto/micro_mall_pay_proto/pay_business"

	"gitee.com/cristiane/micro-mall-users-consumer/model/args"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users-consumer/repository"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
)

func UserStateNoticeConsume(ctx context.Context, body string) error {
	var businessMsg args.CommonBusinessMsg
	var err error
	err = json.Unmarshal(body, &businessMsg)
	if err != nil {
		return err
	}

	var notice args.UserStateNotice
	err = json.Unmarshal(businessMsg.Content, &notice)
	if err != nil {
		return err
	}

	userInfo, err := repository.GetUserInfoByUid(notice.Uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, notice.Uid)
		return err
	}
	if businessMsg.Type == args.UserStateEventTypeMerchantInfo && notice.Extra["operation_type"] == "create" {
		// 为商户初始化账户
		serverName := args.RpcServiceMicroMallPay
		conn, err := util.GetGrpcClient(ctx, serverName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
			return err
		}
		balanceInit := "0.01"
		accountReq := pay_business.CreateAccountRequest{
			Owner:       notice.Extra["merchant_code"],
			AccountType: pay_business.AccountType_Company,
			CoinType:    pay_business.CoinType_CNY,
			Balance:     balanceInit,
		}
		client := pay_business.NewPayBusinessServiceClient(conn)
		accountRsp, err := client.CreateAccount(ctx, &accountReq)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount %v,err: %v, r: %v", serverName, err, json.MarshalToStringNoError(accountReq))
			return err
		}
		if accountRsp.Common.Code != pay_business.RetCode_SUCCESS {
			switch accountRsp.Common.Code {
			case pay_business.RetCode_ACCOUNT_EXIST:
				return nil
			default:
				kelvins.ErrLogger.Errorf(ctx, "CreateAccount  req:%v, rsp: %+v", json.MarshalToStringNoError(accountReq), json.MarshalToStringNoError(accountRsp))
				return fmt.Errorf(accountRsp.Common.Msg)
			}
		}

		// 发送初始个人账户成功邮件
		if userInfo.Email != "" {
			emailNotice := fmt.Sprintf(args.MerchantCreateAccountTemplate, userInfo.UserName, businessMsg.Time, balanceInit, "CNY")
			err = email.SendEmailNotice(ctx, userInfo.Email, kelvins.AppName, emailNotice)
			if err != nil {
				kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
			}
		}
	}

	if userInfo.Email != "" {
		var emailNotice string
		switch businessMsg.Type {
		case args.UserStateEventTypePwdModify:
			emailNotice = fmt.Sprintf(args.UserPwdResetTemplate, userInfo.UserName, businessMsg.Time, notice.Extra["reset_type"])
		case args.UserStateEventTypeLogin:
			emailNotice = fmt.Sprintf(args.UserLoginTemplate, userInfo.UserName, businessMsg.Time, notice.Extra["login_type"])
		case args.UserStateEventTypeAccountCharge:
			emailNotice = fmt.Sprintf(args.UserAccountChargeTemplate, userInfo.UserName, businessMsg.Time, notice.Extra["amount"], notice.Extra["coin"])
		case args.UserStateEventTypeMerchantInfo:
			emailNotice = fmt.Sprintf(args.MerchantApplyInfoTemplate, userInfo.UserName, businessMsg.Time)
		}
		err = email.SendEmailNotice(ctx, userInfo.Email, kelvins.AppName, emailNotice)
		if err != nil {
			kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
		}
	}

	return nil
}

func UserStateNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
