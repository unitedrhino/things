package errors

const OTA_ERROR = 2100000

var (
	OtaRetryStatusError  = NewCodeError(OTA_ERROR+1, "升级状态不允许重新升级")
	OtaCancleStatusError = NewCodeError(OTA_ERROR+2, "升级状态已结束")
	OtaDeviceNumError    = NewCodeError(OTA_ERROR+3, "验证设备数不能超过10个")
)
