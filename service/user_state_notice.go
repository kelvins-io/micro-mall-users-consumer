package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users-consumer/model/args"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users-consumer/repository"
	"gitee.com/cristiane/micro-mall-users-consumer/vars"
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
	if businessMsg.Type != args.UserStateEventTypePwdModify {
		return nil
	}

	var notice args.UserStateNotice
	err = json.Unmarshal(businessMsg.Content, &notice)
	if err != nil {
		return err
	}
	kelvins.ErrLogger.Info(ctx, "UserStateNoticeConsume content: %v", json.MarshalToStringNoError(notice))
	userInfo, err := repository.GetUserNameByUid(notice.Uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, notice.Uid)
		return err
	}

	emailNotice := fmt.Sprintf(args.UserPwdResetTemplate, userInfo.UserName, notice.Time)
	if vars.EmailNoticeSetting != nil && vars.EmailNoticeSetting.Receivers != nil {
		for _, receiver := range vars.EmailNoticeSetting.Receivers {
			err = email.SendEmailNotice(ctx, receiver, kelvins.AppName, emailNotice)
			if err != nil {
				kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
			}
		}
	}

	return nil
}

func UserStateNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
