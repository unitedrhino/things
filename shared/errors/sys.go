package errors

var (
	SysErrorDefault				= NewCodeError(9900001,"其他错误")
	SysErrorTokenExpired 		= NewCodeError(9900002,"token 已经过期")
	SysErrorTokenNotValidYet  	= NewCodeError(9900003,"token还未生效")
	SysErrorTokenMalformed 		= NewCodeError(9900004,"这不是一个token")
	SysErrorTokenInvalid 		= NewCodeError(9900005,"违法的token")
)