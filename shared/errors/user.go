package errors

var (
	ErrorDuplicateUsername 	= NewCodeError(100001,"用户名已经注册")
	ErrorDuplicateMobile   	= NewCodeError(100002,"手机号已经被占用")
	ErrorUsernameUnRegister = NewCodeError(100003,"未注册")
	ErrorPassword 			= NewCodeError(100004,"账号或密码错误")
	ErrorCaptcha 			= NewCodeError(100005,"验证码错误")
	ErrorUidNotCompare 		= NewCodeError(100006,"uid不对应")
	ErrorParameter   		= NewCodeError(100020,"参数错误")
	ErrorSystem   			= NewCodeError(100021,"系统错误")
	ErrorRegisterOne   		= NewCodeError(100022,"注册第一步未成功")
	ErrorDuplicateRegister  = NewCodeError(100022,"重复注册")
)