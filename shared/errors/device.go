package errors

const DEVICE_ERROR = 2000000

var (
	RespParam     = NewCodeError(DEVICE_ERROR+1, "返回参数不对")
	DeviceTimeOut = NewCodeError(DEVICE_ERROR+2, "设备回复超时")
)
