package errors

var (
	SysErrorTokenExpired 		= NewCodeError(9900001,"token 已经过期")
	SysErrorTokenNotValidYet  	= NewCodeError(9900002,"token还未生效")
	SysErrorTokenMalformed 		= NewCodeError(9900003,"这不是一个token")
	SysErrorTokenInvalid 		= NewCodeError(9900004,"违法的token")
)