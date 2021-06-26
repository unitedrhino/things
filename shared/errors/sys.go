package errors

var (
	OK          	 = NewCodeError(0, "成功")
	Default          = NewCodeError(9900001, "其他错误")
	TokenExpired     = NewCodeError(9900002, "token 已经过期")
	TokenNotValidYet = NewCodeError(9900003, "token还未生效")
	TokenMalformed   = NewCodeError(9900004, "这不是一个token")
	TokenInvalid     = NewCodeError(9900005, "违法的token")
	Parameter        = NewCodeError(9900006, "参数错误")
	System           = NewCodeError(9900007, "系统错误")
	Database         = NewCodeError(9900008, "数据库错误")
	NotFind          = NewCodeError(9900009, "未查询到")
	Duplicate        = NewCodeError(99000010, "参数重复")
	SignatureExpired = NewCodeError(99000011, "签名已经过期")
	Permissions      = NewCodeError(99000012, "权限不足")
	Method      	 = NewCodeError(99000013, "method不支持")
	Type      	 = NewCodeError(99000013, "参数的类型不对")
)
