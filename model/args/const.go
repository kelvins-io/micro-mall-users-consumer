package args

type MerchantsMaterialInfo struct {
	Uid          int64
	MaterialId   int64
	RegisterAddr string
	HealthCardNo string
	Identity     int32
	State        int32
	TaxCardNo    string
}

const (
	RpcServiceMicroMallPay = "micro-mall-pay"
	RpcServiceMicroMallUsers = "micro-mall-users"
)

const (
	UserRegisterTemplate      = "用户: +%v-%v, 于%v 注册成功, 状态为: %v"
	UserCreateAccountTemplate = "用户: +%v-%v, 于%v 初始个人账户成功，金额为: %v"
	UserPwdResetTemplate      = "用户: %v, 于%v, 重置密码"
)

const (
	UserStateInit       = 0
	UserStateVerifyIng  = 1
	UserStateVerifyDeny = 2
	UserStateVerifyOk   = 3
)

var UserStateText = map[int]string{
	UserStateInit:       "用户初始化",
	UserStateVerifyIng:  "用户审核中",
	UserStateVerifyDeny: "用户审核拒绝",
	UserStateVerifyOk:   "用户审核通过",
}

const (
	UserStateEventTypeRegister  = 10010
	UserStateEventTypeLogin     = 10011
	UserStateEventTypeLogout    = 10012
	UserStateEventTypePwdModify = 10013
)

type CommonBusinessMsg struct {
	Type int    `json:"type"`
	Tag  string `json:"tag"`
	UUID string `json:"uuid"`
	Msg  string `json:"msg"`
}

type UserRegisterNotice struct {
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	Time        string `json:"time"`
	State       int    `json:"state"`
}

type UserStateNotice struct {
	Uid  int    `json:"uid"`
	Time string `json:"time"`
}
