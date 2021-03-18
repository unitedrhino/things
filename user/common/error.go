package common

import "yl/shared/define"

var (
	ErrorDuplicateUsername 	= define.NewCodeError(100001,"用户名已经注册")
	ErrorDuplicateMobile   	= define.NewCodeError(100002,"手机号已经被占用")
	ErrorUsernameUnRegister = define.NewCodeError(100003,"未注册")
	ErrorPassword 			= define.NewCodeError(100004,"账号或密码错误")
	ErrorCaptcha 			= define.NewCodeError(100005,"验证码错误")
	ErrorUidNotCompare 		= define.NewCodeError(100006,"uid不对应")
	ErrorParameter   		= define.NewCodeError(100020,"参数错误")
	ErrorSystem   			= define.NewCodeError(100021,"系统错误")
	ErrorRegisterOne   		= define.NewCodeError(100022,"注册第一步未成功")
	ErrorDuplicateRegister  = define.NewCodeError(100022,"重复注册")
)