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
	RpcServiceMicroMallPay   = "micro-mall-pay"
	RpcServiceMicroMallUsers = "micro-mall-users"
)

const (
	UserLoginTemplate             = "尊敬的用户【%s】你好，你于：%v 在微商城使用【%s】登录"
	UserRegisterTemplate          = "尊敬的用户【%s】, 于%v 注册成功, 状态为: %v"
	UserCreateAccountTemplate     = "尊敬的用户【%s】你好 , 于%v 初始个人交易账户成功，初始金额为: %v，币种：%v"
	MerchantCreateAccountTemplate = "尊敬的商户【%s】你好 , 于%v 初始公司交易账户成功，初始金额为: %v，币种：%v"
	UserAccountChargeTemplate     = "尊敬的用户【%s】你好 , 于%v 充值个人交易账户成功，充值金额为: %v，币种：%v"
	UserPwdResetTemplate          = "尊敬的用户【%s】你好, 于%v, 通过：%v，重置密码成功"
	MerchantApplyInfoTemplate     = "尊敬的商户【%s】你好，于%v，提交商户认证资料成功"
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
	UserStateEventTypeRegister      = 10010
	UserStateEventTypeLogin         = 10011
	UserStateEventTypeLogout        = 10012
	UserStateEventTypePwdModify     = 10013
	UserStateEventTypeAccountCharge = 10014
	UserStateEventTypeMerchantInfo  = 10015
)

type CommonBusinessMsg struct {
	Type    int    `json:"type"`
	Tag     string `json:"tag"`
	UUID    string `json:"uuid"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

type UserRegisterNotice struct {
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	State       int    `json:"state"`
}

type UserStateNotice struct {
	Uid   int               `json:"uid"`
	Extra map[string]string `json:"extra"`
}
