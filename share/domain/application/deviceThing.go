package application

import (
	"encoding/json"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
)

type ParamValue struct {
	Value any             `json:"value"` //值
	Type  schema.DataType `json:"type"`  //值的类型
}

type StructValue map[string]any

// 固件升级进度上报消息体
type OtaReport struct {
	Device    devices.Core        `json:"device"`
	Timestamp int64               `json:"timestamp,string"` //毫秒时间戳
	Status    msgOta.DeviceStatus `json:"status"`           //设备升级作业状态。1：待确认。 2：待推送。 3：已推送。  4：升级中。 5:升级成功 6: 升级失败. 7:已取消
	Detail    string              `json:"detail"`           //详情
	Step      int64               `json:"step"`             //当前的升级进度  0-100%    -1：升级失败。-2：下载失败。-3：校验失败。-4：烧写失败。
}

// 属性上报消息体
type PropertyReport struct {
	Device     devices.Core `json:"device"`
	Timestamp  int64        `json:"timestamp,string"` //毫秒时间戳
	Identifier string       `json:"identifier"`       //推送属性的标识符
	Param      any          `json:"param"`            //推送属性的参数
}

type PropertyReportV2 struct {
	Device    devices.Core   `json:"device"`
	Timestamp int64          `json:"timestamp,string"` //毫秒时间戳
	Params    map[string]any `json:"params"`           //推送属性的参数,key为属性的id,value为值
}

type Hub struct {
	ProductID   string `json:"productID"`                  // 产品id
	DeviceName  string `json:"deviceName"`                 // 设备名称
	Content     string `json:"content,omitempty"`          // 具体信息
	Topic       string `json:"topic,omitempty"`            // 主题
	Action      string `json:"action,omitempty"`           // 操作类型
	Timestamp   int64  `json:"timestamp,string,omitempty"` // 操作时间,
	RequestID   string `json:"requestID,omitempty"`        // 请求ID
	TraceID     string `json:"traceID,omitempty"`          // 服务器端事务id
	ResultCode  int64  `json:"resultCode,omitempty"`       // 请求结果状态,200为成功
	RespPayload string `json:"respPayload,omitempty"`      //返回的内容
}

// 行为上报消息体
type ActionReport struct {
	Device    devices.Core      `json:"device"`
	MsgToken  string            `json:"msgToken,omitempty"` //调用id
	Timestamp int64             `json:"timestamp,string"`   //毫秒时间戳
	ActionID  string            `json:"actionID,omitempty"` //数据模板中的行为标识符，由开发者自行根据设备的应用场景定义
	Params    map[string]any    `json:"params,omitempty"`   //参数列表
	Code      int64             `json:"code,omitempty"`     //200为成功
	Msg       string            `json:"msg,omitempty"`      //消息
	Dir       schema.ActionDir  `json:"dir"`                //请求方向 up:设备请求云端  down:云端请求设备
	ReqType   deviceMsg.ReqType `json:"reqType"`            //请求类型 req resp
}

// 事件上报消息体
type EventReport struct {
	Device     devices.Core   `json:"device"`
	Timestamp  int64          `json:"timestamp,string"` //毫秒时间戳
	Identifier string         `json:"identifier"`       //标识符
	Type       string         `json:"type" `            //事件类型: 信息:info  告警:alert  故障:fault
	Params     map[string]any `json:"params" `          //事件参数
}

func (c PropertyReport) GenSerial() string {
	return ""
}

func (c EventReport) GenSerial() string {
	return ""
}

//func (c ParamValue) MarshalJSON() ([]byte, error) {
//
//}

func (c *ParamValue) UnmarshalJSON(b []byte) error {
	//自定义json解析
	stu := map[string]any{}
	err := json.Unmarshal(b, &stu)
	if err != nil {
		return err
	}
	ret := parse(stu)
	c.Type = ret.Type
	c.Value = ret.Value
	return nil
}
func parse(stu map[string]any) (ret ParamValue) {
	ret.Type = schema.DataType(cast.ToString(stu["type"]))
	ret.Value = stu["value"]
	switch ret.Type {
	case schema.DataTypeFloat:
		ret.Value = cast.ToFloat64(ret.Value)
		return
	case schema.DataTypeInt, schema.DataTypeEnum, schema.DataTypeTimestamp:
		ret.Value = cast.ToInt64(ret.Value)
		return
	case schema.DataTypeStruct:
		val := ret.Value.(map[string]any)
		structVal := StructValue{}
		for k, v := range val {
			structVal[k] = parse(v.(map[string]any))
		}
		ret.Value = structVal
	case schema.DataTypeArray:

	default:
		return
	}
	return
}
