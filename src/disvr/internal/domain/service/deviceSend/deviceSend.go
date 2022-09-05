package deviceSend

import (
	"encoding/json"
	"fmt"
	"time"
)

type Method = string

const (
	/*
		当设备需要向云端上报设备运行状态的变化时，以通知应用端小程序、App 实时展示或云端业务系统接收设备上报属性数据，物联网开发平台为设备设定了默认的 Topics：
		设备属性上行请求 Topics： $thing/up/property/{ProductID}/{DeviceName}
		设备属性上行响应 Topics： $thing/down/property/{ProductID}/{DeviceName}
	*/
	Report      Method = "report"       //表示设备属性上报
	ReportReply Method = "report_reply" // 表示云端接收设备上报后的响应报文

	/*
		使用数据模板协议的设备，当需要通过云端远程控制设备时，设备需订阅下发 Topics 接收云端指令：
		下发 Topics： $thing/down/property/{ProductID}/{DeviceName}
		响应 Topics： $thing/up/property/{ProductID}/{DeviceName}
	*/
	//todo 尚未支持
	Control      Method = "control"       //表示云端向设备发起控制请求
	ControlReply Method = "control_reply" //表示设备向云端下发的控制指令的请求响应

	/*
		设备从云端接收最新消息使用的 Topics：
		请求 Topics： $thing/up/property/{ProductID}/{DeviceName}
		响应 Topics： $thing/down/property/{ProductID}/{DeviceName}
	*/
	GetStatus      Method = "get_status"       //表示获取设备最新上报的信息
	GetStatusReply Method = "get_status_reply" //表示获取设备最新上报信息的 reply 消息

	/*
		当设备需要向云端上报事件时，如上报设备的故障、告警数据，开发平台为设备设定了默认的 Topics：
		设备事件上行请求 Topics： $thing/up/event/{ProductID}/{DeviceName}
		设备事件上行响应 Topics： $thing/down/event/{ProductID}/{DeviceName}
	*/
	EventPost  Method = "event_post"  //表示事件上报
	EventReply Method = "event_reply" //表示是云端返回设备端的响应

	/*
		当应用通过云端向设备发起某个行为调用时，开发平台为设备行为的处理设定了默认的 Topics：
		应用调用设备行为 Topics： $thing/down/action/{ProductID}/{DeviceName}
		设备响应行为执行结果 Topics： $thing/up/action/{ProductID}/{DeviceName}
	*/
	Action      Method = "action"       //表示是调用设备的某个行为
	ActionReply Method = "action_reply" //表示是设备端执行完指定的行为向云端回复的响应

	/*
		小程序或 App 展示设备详细信息时，一般会展示设备的 MAC 地址、IMEI 号、时区等基础信息。设备信息上报使用的 Topics：
		上行请求 Topics： $thing/up/property/{ProductID}/{DeviceName}
		下行响应 Topics： $thing/down/property/{ProductID}/{DeviceName}
	*/
	// todo 尚未支持
	ReportInfo      Method = "report_info"       //表示设备基础信息上报
	ReportInfoReply Method = "report_info_reply" //表示云端接收设备上报后的响应报文

	/*
		当用户在小程序或App中删除已绑定成功的设备，平台会下发用户删除设备的通知到设备，设备接收后可根据业务需求自行处理。如网关类设备接收到子设备被删除。
		下发用户删除设备 Topics： $thing/down/service/{ProductID}/{DeviceName}
	*/
	UnbindDevice Method = "unbind_device" // 表示是用户在小程序或 App 中删除或解绑某个设备

)

func GetMethod(method Method) Method {
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
	default:
		panic(method)
	}
}

// GenThingDeviceRespData 生成物模型设备请求的回复包
func GenThingDeviceRespData(Method, ClientToken string, topics []string, err error,
	data map[string]any) (topic string, payload []byte) {
	respMethod := GetMethod(Method)
	respTopic := fmt.Sprintf("%s/down/%s/%s/%s",
		topics[0], topics[2], topics[3], topics[4])
	respPayload, _ := json.Marshal(DeviceResp{
		Method:      respMethod,
		ClientToken: ClientToken,
		Data:        data,
		Timestamp:   time.Now().UnixMilli(),
	}.AddStatus(err))
	return respTopic, respPayload
}
