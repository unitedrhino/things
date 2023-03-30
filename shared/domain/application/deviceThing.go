package application

import (
	"encoding/json"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/spf13/cast"
)

type ParamValue struct {
	Value any             `json:"value"` //值
	Type  schema.DataType `json:"type"`  //值的类型
}

type StructValue map[string]ParamValue

//属性上报消息体
type PropertyReport struct {
	Device     devices.Core `json:"device"`
	Timestamp  int64        `json:"timestamp,string"` //毫秒时间戳
	Identifier string       `json:"identifier"`       //推送属性的标识符
	Param      ParamValue   `json:"param"`            //推送属性的参数
}

type EventReport struct {
	Device     devices.Core          `json:"device"`
	Timestamp  int64                 `json:"timestamp,string"` //毫秒时间戳
	Identifier string                `json:"identifier"`       //标识符
	Type       string                `json:"type" `            //事件类型: 信息:info  告警:alert  故障:fault
	Params     map[string]ParamValue `json:"params" `          //事件参数
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
