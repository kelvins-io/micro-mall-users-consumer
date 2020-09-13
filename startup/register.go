package startup

import "gitee.com/cristiane/micro-mall-users-consumer/service"

const (
	TaskNameUserRegisterNotice    = "task_user_register_notice"
	TaskNameUserRegisterNoticeErr = "task_user_register_notice_err"

	TaskNameUserStateNotice    = "task_user_state_notice"
	TaskNameUserStateNoticeErr = "task_user_state_notice_err"
)

func GetNamedTaskFuncs() map[string]interface{} {

	var taskRegister = map[string]interface{}{
		TaskNameUserRegisterNotice:    service.UserRegisterNoticeConsume,
		TaskNameUserRegisterNoticeErr: service.UserRegisterNoticeConsumeErr,
		TaskNameUserStateNotice:       service.UserStateNoticeConsume,
		TaskNameUserStateNoticeErr:    service.UserStateNoticeConsumeErr,
	}
	return taskRegister
}
