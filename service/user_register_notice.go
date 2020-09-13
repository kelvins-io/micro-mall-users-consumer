package service

import (
	"context"
	"gitee.com/kelvins-io/kelvins"
)

func UserRegisterNoticeConsume(ctx context.Context, body string) error {
	kelvins.ErrLogger.Info(ctx, "UserRegisterNoticeConsume body: %v", body)
	return nil
}

func UserRegisterNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
