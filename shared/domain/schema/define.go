package schema

//数据类型
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

//物模型标签
type Tag int64

const (
	TagCustom   Tag = 1 //自定义
	TagOptional Tag = 2 //可选
	TagRequired Tag = 3 //必选 必选不可删除
)

//物模型功能类型 1:property属性 2:event事件 3:action行为
type AffordanceType int64

const (
	//物模型功能类型：1-property 属性
	AffordanceTypeProperty AffordanceType = 1
	//物模型功能类型：2-event 事件
	AffordanceTypeEvent AffordanceType = 2
	//物模型功能类型：3-action 行为
	AffordanceTypeAction AffordanceType = 3
)

func (m AffordanceType) String() string {
	switch m {
	case AffordanceTypeProperty:
		return "property"
	case AffordanceTypeEvent:
		return "event"
	case AffordanceTypeAction:
		return "action"

	}
	return ""
}

//属性读写类型: r(只读) rw(可读可写)
type PropertyMode string

const (
	PropertyModeR  PropertyMode = "r"
	PropertyModeRW PropertyMode = "rw"
)

//事件类型: 信息:info  告警alert  故障:fault
type EventType string

const (
	EventTypeInfo  EventType = "info"
	EventTypeAlert EventType = "alert"
	EventTypeFault EventType = "fault"
)
