package errors

const OtaError = 2100000

var (
	OtaRetryStatusError  = NewCodeError(OtaError+1, "升级状态不允许重新升级")
	OtaCancleStatusError = NewCodeError(OtaError+2, "升级状态已结束")
)
