package deviceMsg

type Method = string

const (
	/*
		当设备需要向云端上报设备运行状态的变化时，以通知应用端小程序、App 实时展示或云端业务系统接收设备上报属性数据，物联网开发平台为设备设定了默认的 Topic：
		设备属性上行请求 Topic： $thing/up/property/{ProductID}/{DeviceNames}
		设备属性下行响应 Topic： $thing/down/property/{ProductID}/{DeviceNames}
	*/
	Report      Method = "report"      //表示设备属性上报
	ReportReply Method = "reportReply" // 表示云端接收设备上报后的响应报文

	/*
		需要主动请求设备上报的时候需要用以下方式上报 Topic：
		设备属性下行请求 Topic： $thing/up/property/{ProductID}/{DeviceNames}
		设备属性上行响应 Topic： $thing/down/property/{ProductID}/{DeviceNames}
	*/
	GetReport      Method = "getReport"      //表示云端请求设备获取上报消息
	GetReportReply Method = "getReportReply" // 表示设备属性上报

	/*
		使用数据模板协议的设备，当需要通过云端远程控制设备时，设备需订阅下发 Topic 接收云端指令：
		下发 Topic： $thing/down/property/{ProductID}/{DeviceNames}
		响应 Topic： $thing/up/property/{ProductID}/{DeviceNames}
	*/
	Control      Method = "control"      //表示云端向设备发起控制请求
	ControlReply Method = "controlReply" //表示设备向云端下发的控制指令的请求响应（设备回复的 云端下发控制指令 的处理结果）

	/*
		设备从云端接收最新消息使用的 Topic：
		请求 Topic： $thing/up/property/{ProductID}/{DeviceNames}
		响应 Topic： $thing/down/property/{ProductID}/{DeviceNames}
	*/
	GetStatus      Method = "getStatus"      //表示获取设备最新上报的信息（设备请求获取 云端记录的最新设备信息）
	GetStatusReply Method = "getStatusReply" //表示获取设备最新上报信息的 reply 消息

	/*
		当设备需要向云端上报事件时，如上报设备的故障、告警数据，开发平台为设备设定了默认的 Topic：
		设备事件上行请求 Topic： $thing/up/event/{ProductID}/{DeviceNames}
		设备事件上行响应 Topic： $thing/down/event/{ProductID}/{DeviceNames}
	*/
	EventPost  Method = "eventPost"  //表示事件上报
	EventReply Method = "eventReply" //表示是云端返回设备端的响应

	/*
		当应用通过云端向设备发起某个行为调用时，开发平台为设备行为的处理设定了默认的 Topic：
		应用调用设备行为 Topic： $thing/down/action/{ProductID}/{DeviceNames}
		设备响应行为执行结果 Topic： $thing/up/action/{ProductID}/{DeviceNames}
	*/
	Action      Method = "action"      //表示是调用设备的某个行为
	ActionReply Method = "actionReply" //表示是设备端执行完指定的行为向云端回复的响应

	/*
		小程序或 App 展示设备详细信息时，一般会展示设备的 MAC 地址、IMEI 号、时区等基础信息。设备信息上报使用的 Topic：
		上行请求 Topic： $thing/up/property/{ProductID}/{DeviceNames}
		下行响应 Topic： $thing/down/property/{ProductID}/{DeviceNames}
	*/
	// todo 尚未支持
	ReportInfo      Method = "reportInfo"      //表示设备基础信息上报
	ReportInfoReply Method = "reportInfoReply" //表示云端接收设备上报后的响应报文

	/*
		拓扑关系管理
		网关类型的设备，可通过与云端的数据通信，对其下的子设备进行绑定与解绑操作。实现此类功能需利用如下两个 Topic：
		数据上行 Topic（用于发布）：$gateway/operation/${productid}/${devicename}
		数据下行 Topic（用于订阅）：$gateway/operation/${productid}/${devicename}
	*/
	Bind               Method = "bind"               //绑定设备
	Unbind             Method = "unbind"             //解绑设备
	DescribeSubDevices Method = "describeSubDevices" //查询拓扑关系
	Change             Method = "change"             //拓扑关系变化
	Register           Method = "register"           //注册新设备
	/*
		数据上行 Topic（用于发布）：$gateway/status/${productid}/${devicename}
		数据下行 Topic（用于订阅）：$gateway/status/${productid}/${devicename}
	*/
	Online  Method = "online"  //代理子设备上线
	Offline Method = "offline" //代理子设备下线

	/*
		数据上行 Topic（用于发布）：$config/up/get/${productid}/${devicename}
		数据下行 Topic（用于订阅）：$config/down/get/${productid}/${devicename}
	*/
	RemoteConfigReply Method = "reply" //表示设备请求配置

	/*
		ntp时间同步
		设备ntp上行请求 Topic（用于发布）：$ext/up/ntp/${productid}/${devicename}
		设备ntp下行响应 Topic（用于订阅）：$ext/down/ntp/${productid}/${devicename}
	*/
	GetNtp      Method = "getNtp"      //表示设备ntp请求
	GetNtpReply Method = "getNtpReply" //表示云端ntp返回
)

func GetRespMethod(method Method) Method {
	switch method {
	case ReportInfo:
		return ReportInfoReply
	case Action:
		return ActionReply
	case EventPost:
		return EventReply
	case GetStatus:
		return GetStatusReply
	case Control:
		return ControlReply
	case Report:
		return ReportReply
	case GetNtp:
		return GetNtpReply
	default: //不在里面的方法直接返回方法名即可
		return method
	}
}
