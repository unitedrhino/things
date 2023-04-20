package msgThing

const (
	TypeProperty = "property" //事件-操作类型：属性或信息上报（Topic：$thing/up/property/{ProductID}/{DeviceName}）
	TypeEvent    = "event"    //事件-操作类型：事件上报（Topic：$thing/up/event/{ProductID}/{DeviceName}）
	TypeAction   = "action"   //事件-操作类型：行为上报（Topic：$thing/up/action/{ProductID}/{DeviceName}）

	EventReport = "report" //设备上报的信息
	EventInfo   = "info"   //信息
	EventAlert  = "alert"  //告警
	EventFault  = "fault"  //故障
)
