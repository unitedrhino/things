package errors

const USER_ERROR = 1000000

var (
	DuplicateUsername  = NewCodeError(USER_ERROR+1, "用户名已经注册")
	DuplicateMobile    = NewCodeError(USER_ERROR+2, "手机号已经被占用")
	UnRegister         = NewCodeError(USER_ERROR+3, "未注册")
	Password           = NewCodeError(USER_ERROR+4, "账号或密码错误")
	Captcha            = NewCodeError(USER_ERROR+5, "验证码错误")
	UidNotRight        = NewCodeError(USER_ERROR+6, "uid不对")
	NotLogin           = NewCodeError(USER_ERROR+7, "尚未登录")
	RegisterOne        = NewCodeError(USER_ERROR+22, "注册第一步未成功")
	DuplicateRegister  = NewCodeError(USER_ERROR+23, "重复注册")
	NeedUserName       = NewCodeError(USER_ERROR+24, "需要填入用户名")
	PasswordLevel      = NewCodeError(USER_ERROR+25, "密码强度不够")
	GetInfoPartFailure = NewCodeError(USER_ERROR+26, "获取用户信息有失败")
	UsernameFormatErr  = NewCodeError(USER_ERROR+27, "账密方式时，账号必须以大小写字母开头，且账号只能包含大小写字母，数字，下划线和减号。 长度为6到20位之间")
	AccountForbidden   = NewCodeError(USER_ERROR+28, "用户账号冻结")
	IpForbidden        = NewCodeError(USER_ERROR+29, "ip冻结")
	UseCaptcha         = NewCodeError(USER_ERROR+30, "连续密码错误触发验证码")
)
