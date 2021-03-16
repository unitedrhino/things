package logic

import "yl/shared/define"

var (
	errorDuplicateUsername 	= define.NewCodeError(100001,"用户名已经注册")
	errorDuplicateMobile   	= define.NewCodeError(100002,"手机号已经被占用")
	errorUsernameUnRegister = define.NewCodeError(100003,"未注册")
	errorPassword 			= define.NewCodeError(100004,"账号或密码错误")
	errorCaptcha 			= define.NewCodeError(100005,"验证码错误")
	errorParameter   		= define.NewCodeError(100020,"参数错误")
	errorSystem   			= define.NewCodeError(100021,"系统错误")
	errorRegisterOne   		= define.NewCodeError(100022,"注册第一步未成功")
	errorDuplicateRegister  = define.NewCodeError(100022,"重复注册")
)