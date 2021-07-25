package startup

import (
	"gitee.com/cristiane/micro-mall-users-consumer/vars"
	"gitee.com/kelvins-io/kelvins/config"
)

const (
	SectionEmailConfig             = "email-config"
	//SectionQueueUserRegisterNotice = "queue-user-register-notice"
	//SectionQueueUserStateNotice    = "queue-user-state-notice"
	EmailNotice                    = "email-notice"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 邮件通知
	vars.EmailNoticeSetting = new(vars.EmailNoticeSettingS)
	config.MapConfig(EmailNotice, vars.EmailNoticeSetting)
	return nil
}
