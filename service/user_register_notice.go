package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users-consumer/model/args"
	"gitee.com/cristiane/micro-mall-users-consumer/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users-consumer/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
)

const (
	templateUserRegister = "用户: +%v-%v, 于%v 注册成功, 状态为: %v"
)

func UserRegisterNoticeConsume(ctx context.Context, body string) error {
	var businessMsg args.CommonBusinessMsg
	var err error
	err = json.Unmarshal(body, &businessMsg)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "第1次序列化失败 err: %v", err)
		return err
	}
	if businessMsg.Type != args.UserStateEventTypeRegister {
		return nil
	}
	kelvins.ErrLogger.Info(ctx, "businessMsg: %+v", businessMsg)
	var notice args.UserRegisterNotice
	err = json.Unmarshal(businessMsg.Msg, &notice)
	if err != nil {
		kelvins.ErrLogger.Info(ctx, "第2次序列化失败 err: %v", err)
		return err
	}

	kelvins.ErrLogger.Info(ctx, "UserRegisterNoticeConsume content: %v", notice)
	kelvins.ErrLogger.Info(ctx, "开始发送邮件")
	err = email.SendEmailNotice(ctx, "565608463@qq.com",
		vars.AppName, fmt.Sprintf(templateUserRegister, notice.CountryCode, notice.Phone, notice.Time, notice.State))

	if err != nil {
		kelvins.ErrLogger.Info(ctx, "邮件发送失败")
		return err
	}
	kelvins.ErrLogger.Info(ctx, "邮件发送成功")

	return nil
}

func UserRegisterNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
