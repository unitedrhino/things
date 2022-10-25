package errors

const LOG_ERROR = 3000000

var (
	SDKLOG_MissParam  = NewCodeError(LOG_ERROR+1, "topic参数缺失")
	SDKLOG_ErrorLevel = NewCodeError(LOG_ERROR+1, "loglevel参数错误")
)
