package errors

var (
	DuplicateUsername 	= NewCodeError(100001,"用户名已经注册")
	DuplicateMobile   	= NewCodeError(100002,"手机号已经被占用")
	UnRegister  		= NewCodeError(100003,"未注册")
	Password 			= NewCodeError(100004,"账号或密码错误")
	Captcha 			= NewCodeError(100005,"验证码错误")
	UidNotCompare 		= NewCodeError(100006,"uid不对应")
	RegisterOne   		= NewCodeError(100022,"注册第一步未成功")
	DuplicateRegister   = NewCodeError(100023,"重复注册")
	NeedUserName   		= NewCodeError(100024,"需要填入用户名")
	PasswordLevel   	= NewCodeError(100025,"密码强度不够")
)