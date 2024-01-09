package errors

const SysError = 100000

var (
	OK               = NewCodeError(200, "成功")
	Default          = NewCodeError(SysError+1, "其他错误")
	TokenExpired     = NewCodeError(SysError+2, "token已经过期")
	TokenNotValidYet = NewCodeError(SysError+3, "token还未生效")
	TokenMalformed   = NewCodeError(SysError+4, "token格式错误")
	TokenInvalid     = NewCodeError(SysError+5, "违法的token")
	Parameter        = NewCodeError(SysError+6, "参数错误")
	System           = NewCodeError(SysError+7, "系统错误")
	Database         = NewCodeError(SysError+8, "数据库错误")
	NotFind          = NewCodeError(SysError+9, "未查询到")
	Duplicate        = NewCodeError(SysError+10, "参数重复")
	SignatureExpired = NewCodeError(SysError+11, "签名已经过期")
	Permissions      = NewCodeError(SysError+12, "权限不足")
	Method           = NewCodeError(SysError+13, "method不支持")
	Type             = NewCodeError(SysError+14, "参数的类型不对")
	OutRange         = NewCodeError(SysError+15, "参数的值超出范围")
	TimeOut          = NewCodeError(SysError+16, "等待超时")
	Server           = NewCodeError(SysError+17, "本实例处理不了该信息")
	NotRealize       = NewCodeError(SysError+18, "尚未实现")
	NotEmpty         = NewCodeError(SysError+19, "不为空")
	Panic            = NewCodeError(SysError+20, "系统异常，请联系开发者")
	NotEnable        = NewCodeError(SysError+21, "未启用")
	Company          = NewCodeError(SysError+22, "该功能是企业版功能")
	Script           = NewCodeError(SysError+23, "脚本执行失败")
	OnGoing          = NewCodeError(SysError+24, "正在执行中")     //事务分布式事务中如果返回该错误码,分布式事务会定时重试
	Failure          = NewCodeError(SysError+25, "执行失败,需要回滚") //事务分布式事务中如果返回该错误码,分布式事务会进行回滚

)
