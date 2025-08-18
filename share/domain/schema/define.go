package schema

// 数据类型
type DataType string

const (
	DataTypeBool      DataType = "bool"
	DataTypeInt       DataType = "int"
	DataTypeString    DataType = "string"
	DataTypeStruct    DataType = "struct"
	DataTypeFloat     DataType = "float"
	DataTypeTimestamp DataType = "timestamp"
	DataTypeArray     DataType = "array"
	DataTypeEnum      DataType = "enum"
)

type ParamType int64

const (
	//请求参数的类型：设备属性上报参数
	ParamProperty ParamType = iota + 1
	//请求参数的类型：设备行为调用的下行参数
	ParamActionInput
	//请求参数的类型：设备行为调用上行的回复参数
	ParamActionOutput
	//请求参数的类型：设备事件上报参数
	ParamEvent
)

// 物模型标签
type Tag = int64

const (
	TagCustom         Tag = 1 //自定义
	TagOptional       Tag = 2 //可选
	TagRequired       Tag = 3 //必选 必选不可删除
	TagDeviceCustom   Tag = 4 //设备自定义物模型
	TagDeviceOptional Tag = 5 //设备自选物模型
)

type RecordMode = int64

const (
	RecordModeAll  RecordMode = 1 //全部记录
	RecordModeAuto RecordMode = 2 //只记录差异值
	RecordModeNone RecordMode = 3 //不记录
)

// 属性读写类型: r(只读) rw(可读可写)
type PropertyMode string

const (
	PropertyModeR  PropertyMode = "r"
	PropertyModeRW PropertyMode = "rw"
)

// 事件类型: 信息:info  告警alert  故障:fault
type EventType = string

const (
	EventTypeInfo  EventType = "info"
	EventTypeAlert EventType = "alert"
	EventTypeFault EventType = "fault"
)

// 行为的执行方向
type ActionDir = string

const (
	ActionDirUp   ActionDir = "up"   //向上调用
	ActionDirDown ActionDir = "down" //向下调用
)
