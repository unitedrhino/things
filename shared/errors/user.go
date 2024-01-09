package errors

const UserError = 1000000

var (
	DuplicateUsername    = NewCodeError(UserError+1, "用户名已经注册")
	DuplicateMobile      = NewCodeError(UserError+2, "手机号已经被占用")
	UnRegister           = NewCodeError(UserError+3, "未注册")
	Password             = NewCodeError(UserError+4, "账号或密码错误")
	Captcha              = NewCodeError(UserError+5, "验证码错误")
	UidNotRight          = NewCodeError(UserError+6, "uid不对")
	NotLogin             = NewCodeError(UserError+7, "尚未登录")
	RegisterOne          = NewCodeError(UserError+22, "注册第一步未成功")
	DuplicateRegister    = NewCodeError(UserError+23, "重复注册")
	NeedUserName         = NewCodeError(UserError+24, "需要填入用户名")
	PasswordLevel        = NewCodeError(UserError+25, "密码强度不够")
	GetInfoPartFailure   = NewCodeError(UserError+26, "获取用户信息有失败")
	UsernameFormatErr    = NewCodeError(UserError+27, "账号必须以大小写字母开头，且账号只能包含大小写字母，数字，下划线和减号。 长度为6到20位之间")
	AccountOrIpForbidden = NewCodeError(UserError+28, "密码输入错误过多，账号冻结")
	UseCaptcha           = NewCodeError(UserError+29, "账号或密码错误")
)
