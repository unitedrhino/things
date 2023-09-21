package errors

const OTA_ERROR = 2100000

var (
	OtaRetryStatusError  = NewCodeError(OTA_ERROR+1, "升级状态不允许重新升级")
	OtaCancleStatusError = NewCodeError(OTA_ERROR+2, "升级状态已结束")
)
