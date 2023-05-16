package errors

const DEVICE_ERROR = 2000000

var (
	RespParam     = NewCodeError(DEVICE_ERROR+1, "返回参数不对")
	DeviceTimeOut = NewCodeError(DEVICE_ERROR+2, "设备回复超时")
	NotOnline     = NewCodeError(DEVICE_ERROR+3, "设备不在线")
	DeviceError   = NewCodeError(DEVICE_ERROR+4, "设备回复错误")
)
