package product

type DeviceSchemaMode = int64

const (
	DeviceSchemaModeUser             = 1 //用户自己创建
	DeviceSchemaModeAutoCreate       = 2 //设备自动创建
	DeviceSchemaModeReportAutoCreate = 3 //设备自动创建及上报无定义自动创建
)
