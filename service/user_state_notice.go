package service

import (
	"context"
	"gitee.com/kelvins-io/kelvins"
)

func UserStateNoticeConsume(ctx context.Context, body string) error {
	kelvins.ErrLogger.Info(ctx, "UserStateNoticeConsume body: %v", body)
	return nil
}

func UserStateNoticeConsumeErr(ctx context.Context, err, body string) error {
	return nil
}
