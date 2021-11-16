package errors

const USER_ERROR = 1000000

var (
	DuplicateUsername  = NewCodeError(USER_ERROR+1, "用户名已经注册")
	DuplicateMobile    = NewCodeError(USER_ERROR+2, "手机号已经被占用")
	UnRegister         = NewCodeError(USER_ERROR+3, "未注册")
	Password           = NewCodeError(USER_ERROR+4, "账号或密码错误")
	Captcha            = NewCodeError(USER_ERROR+5, "验证码错误")
	UidNotRight        = NewCodeError(USER_ERROR+6, "uid不对")
	RegisterOne        = NewCodeError(USER_ERROR+22, "注册第一步未成功")
	DuplicateRegister  = NewCodeError(USER_ERROR+23, "重复注册")
	NeedUserName       = NewCodeError(USER_ERROR+24, "需要填入用户名")
	PasswordLevel      = NewCodeError(USER_ERROR+25, "密码强度不够")
	GetInfoPartFailure = NewCodeError(USER_ERROR+26, "获取用户信息有失败")
)
