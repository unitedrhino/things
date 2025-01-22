package product

type DeviceSchemaMode = int64

const (
	DeviceSchemaModeUser                     = 1 //用户自己创建
	DeviceSchemaModeAutoCreate               = 2 //设备自动创建
	DeviceSchemaModeReportAutoCreate         = 3 //设备自动创建及上报无定义自动创建
	DeviceSchemaModeReportAutoCreateUseFloat = 4 //设备自动创建及上报无定义自动创建,数字类型只使用浮点
)

type BindLevel = int64

const (
	BindLeveHard1   = 1 //强绑定(默认,只有用户解绑之后才能绑定)
	BindLeveMiddle2 = 2 //中绑定(可以通过token强制绑定设备)
	BindLeveWeak3   = 3 //弱绑定(app可以内部解绑被绑定的设备)
)
