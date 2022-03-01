package device

import (
	"github.com/go-things/things/shared/errors"
	"time"
)

type (
	DeviceReq struct {
		Method      string                 `json:"method"`             //操作方法
		ClientToken string                 `json:"clientToken"`        //方便排查随机数
		Params      map[string]interface{} `json:"params,omitempty"`   //参数列表
		Version     string                 `json:"version,omitempty"`  //协议版本，默认为1.0。
		EventID     string                 `json:"eventId,omitempty"`  //事件的 Id，在数据模板事件中定义。
		ActionID    string                 `json:"actionId,omitempty"` //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
		Timestamp   int64                  `json:"timestamp,omitempty"`
		Showmeta    int64                  `json:"showmeta,omitempty"` //标识回复消息是否带 metadata，缺省为0表示不返回 metadata
		Type        string                 `json:"type,omitempty"`     //	表示获取什么类型的信息。report:表示设备上报的信息 info:信息 alert:告警 fault:故障
		Code        int64                  `json:"code,omitempty"`     //状态码
		Status      string                 `json:"status,omitempty"`   //返回信息
	}
	DeviceResp struct {
		Method      string                 `json:"method"`             //操作方法
		ClientToken string                 `json:"clientToken"`        //方便排查随机数
		Version     string                 `json:"version,omitempty"`  //协议版本，默认为1.0。
		Code        int64                  `json:"code"`               //状态码
		Status      string                 `json:"status,omitempty"`   //返回信息
		Type        string                 `json:"type,omitempty"`     //	表示什么类型的信息。report:表示设备上报的信息
		Data        map[string]interface{} `json:"data,omitempty"`     //返回具体设备上报的最新数据内容
		Response    map[string]interface{} `json:"response,omitempty"` //设备行为中定义的返回参数，设备行为执行成功后，向云端返回执行结果
	}
)

func (d DeviceResp) AddStatus(err error) DeviceResp {
	e := errors.Fmt(err)
	d.Code = e.Code
	d.Status = e.GetDetailMsg()
	return d
}

func (d DeviceReq) AddStatus(err error) DeviceReq {
	e := errors.Fmt(err)
	d.Code = e.Code
	d.Status = e.GetDetailMsg()
	return d
}

func (d *DeviceReq) GetTimeStamp(defaultTime time.Time) time.Time {
	if d.Timestamp == 0 {
		return defaultTime
	}
	return time.UnixMilli(d.Timestamp)
}

type METHOD = string

const (
	/*
		当设备需要向云端上报设备运行状态的变化时，以通知应用端小程序、App 实时展示或云端业务系统接收设备上报属性数据，物联网开发平台为设备设定了默认的 Topic：
		设备属性上行请求 Topic： $thing/up/property/{ProductID}/{DeviceName}
		设备属性上行响应 Topic： $thing/down/property/{ProductID}/{DeviceName}
	*/
	REPORT       METHOD = "report"       //表示设备属性上报
	REPORT_REPLY METHOD = "report_reply" // 表示云端接收设备上报后的响应报文

	/*
		使用数据模板协议的设备，当需要通过云端远程控制设备时，设备需订阅下发 Topic 接收云端指令：
		下发 Topic： $thing/down/property/{ProductID}/{DeviceName}
		响应 Topic： $thing/up/property/{ProductID}/{DeviceName}
	*/
	//todo 尚未支持
	CONTROL       METHOD = "control"       //表示云端向设备发起控制请求
	CONTROL_REPLY METHOD = "control_reply" //表示设备向云端下发的控制指令的请求响应

	/*
		设备从云端接收最新消息使用的 Topic：
		请求 Topic： $thing/up/property/{ProductID}/{DeviceName}
		响应 Topic： $thing/down/property/{ProductID}/{DeviceName}
	*/
	GET_STATUS       METHOD = "get_status"       //表示获取设备最新上报的信息
	GET_STATUS_REPLY METHOD = "get_status_reply" //表示获取设备最新上报信息的 reply 消息

	/*
		当设备需要向云端上报事件时，如上报设备的故障、告警数据，开发平台为设备设定了默认的 Topic：
		设备事件上行请求 Topic： $thing/up/event/{ProductID}/{DeviceName}
		设备事件上行响应 Topic： $thing/down/event/{ProductID}/{DeviceName}
	*/
	EVENT_POST  METHOD = "event_post"  //表示事件上报
	EVENT_REPLY METHOD = "event_reply" //表示是云端返回设备端的响应

	/*
		当应用通过云端向设备发起某个行为调用时，开发平台为设备行为的处理设定了默认的 Topic：
		应用调用设备行为 Topic： $thing/down/action/{ProductID}/{DeviceName}
		设备响应行为执行结果 Topic： $thing/up/action/{ProductID}/{DeviceName}
	*/
	ACTION       METHOD = "action"       //表示是调用设备的某个行为
	ACTION_REPLY METHOD = "action_reply" //表示是设备端执行完指定的行为向云端回复的响应

	/*
		小程序或 App 展示设备详细信息时，一般会展示设备的 MAC 地址、IMEI 号、时区等基础信息。设备信息上报使用的 Topic：
		上行请求 Topic： $thing/up/property/{ProductID}/{DeviceName}
		下行响应 Topic： $thing/down/property/{ProductID}/{DeviceName}
	*/
	// todo 尚未支持
	REPORT_INFO       METHOD = "report_info"       //表示设备基础信息上报
	REPORT_INFO_REPLY METHOD = "report_info_reply" //表示云端接收设备上报后的响应报文

	/*
		当用户在小程序或App中删除已绑定成功的设备，平台会下发用户删除设备的通知到设备，设备接收后可根据业务需求自行处理。如网关类设备接收到子设备被删除。
		下发用户删除设备 Topic： $thing/down/service/{ProductID}/{DeviceName}
	*/
	UNBIND_DEVICE METHOD = "unbind_device" // 表示是用户在小程序或 App 中删除或解绑某个设备

)

func GetMethod(method METHOD) METHOD {
	switch method {
	case REPORT_INFO:
		return REPORT_INFO_REPLY
	case ACTION:
		return ACTION_REPLY
	case EVENT_POST:
		return EVENT_REPLY
	case GET_STATUS:
		return GET_STATUS_REPLY
	case CONTROL:
		return CONTROL_REPLY
	case REPORT:
		return REPORT_REPLY
	default:
		panic(method)
	}
}
