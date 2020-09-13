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

const (
	templateUserPwdReset = "用户: %v, 于%v, 重置密码"
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
	err = json.Unmarshal(businessMsg.Msg, &notice)
	if err != nil {
		return err
	}
	kelvins.ErrLogger.Info(ctx, "UserStateNoticeConsume content: %v", notice)
	userInfo, err := repository.GetUserByUid(notice.Uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, notice.Uid)
		return err
	}
	err = email.SendEmailNotice(ctx, "565608463@qq.com",
		vars.AppName, fmt.Sprintf(templateUserPwdReset, userInfo.UserName, notice.Time))
	if err != nil {
		return err
	}

	return nil
}

func UserStateNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
