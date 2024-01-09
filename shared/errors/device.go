package errors

const DeviceError = 2000000

var (
	RespParam     = NewCodeError(DeviceError+1, "返回参数不对")
	DeviceTimeOut = NewCodeError(DeviceError+2, "设备回复超时")
	NotOnline     = NewCodeError(DeviceError+3, "设备不在线")
	DeviceResp    = NewCodeError(DeviceError+4, "设备回复错误")
)
