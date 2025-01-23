package msgThing

import (
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"math"
)

const (
	//是否校验参数范围
	validateDataRange = true
)

type TimeParam struct {
	Timestamp int64            `json:"timestamp,omitempty"` //毫秒时间戳
	EventID   string           `json:"eventID,omitempty"`   //事件的 Id，在数据模板事件中定义。
	Type      schema.EventType `json:"type,omitempty"`      //事件类型: 信息:info  告警alert  故障:fault
	Params    map[string]Param `json:"params"`
}

type Param struct {
	Identifier string              `json:"identifier"` //标识符
	Name       string              `json:"name"`       //功能名称
	Desc       string              `json:"gesc"`       //描述
	Mode       schema.PropertyMode `json:"mode"`       //读写乐行:rw(可读可写) r(只读)
	Required   bool                `json:"required"`   //是否必须
	Type       schema.EventType    `json:"type"`       //事件类型: 信息:info  告警alert  故障:fault
	Define     *schema.Define
	Value      any `json:"value"` //数据定义
}

func (tp *Param) SetByDefine(d *schema.Define, val any) (err error) {
	tp.Define = d
	tp.Value, err = GetVal(d, val)
	return err
}

func ToVal(tp map[string]Param) (map[string]any, error) {
	ret := make(map[string]any, len(tp))
	var err error
	for k, v := range tp {
		ret[k], err = v.ToVal()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func IsParamValEq(d *schema.Define, v1 any, v2 any) bool {
	if d.Type == schema.DataTypeArray {
		return IsParamValEq(d.ArrayInfo, v1, v2)
	}
	var err error
	v1, err = GetVal(d, v1)
	if err != nil {
		return false
	}
	v2, err = GetVal(d, v2)
	if err != nil {
		return false
	}
	switch d.Type {
	case schema.DataTypeBool:
		return cast.ToBool(v1) == cast.ToBool(v2)
	case schema.DataTypeInt, schema.DataTypeEnum, schema.DataTypeTimestamp:
		return cast.ToInt64(v1) == cast.ToInt64(v2)
	case schema.DataTypeFloat:
		return cast.ToFloat64(v1) == cast.ToFloat64(v2)
	case schema.DataTypeString:
		return cast.ToString(v1) == cast.ToString(v2)
	case schema.DataTypeStruct:
		v1m, ok := v1.(map[string]any)
		if !ok {
			return false
		}
		v2m, ok := v2.(map[string]any)
		if !ok {
			return false
		}
		if len(v1m) != len(v2m) {
			return false
		}
		for k, v := range v1m {
			if !IsParamValEq(&d.Spec[k].DataType, v, v2m[k]) {
				return false
			}
		}
		return true
	}
	return false
}

func getParamVal(def *schema.Define, value any) (any, error) {
	var err error
	switch def.Type {
	case schema.DataTypeBool:
		if cast.ToBool(value) {
			return 1, nil
		}
		return 0, nil
	case schema.DataTypeStruct:
		v, ok := value.(map[string]Param)
		if ok == false {
			return nil, errors.Parameter.AddMsgf("struct Param is not find")
		}
		val := make(map[string]any, len(v)+1)
		for _, tp := range v {
			val[tp.Identifier], err = tp.ToVal()
			if err != nil {
				return nil, err
			}
		}
		return val, nil
	case schema.DataTypeArray:
		array, ok := value.([]any)
		if ok == false {
			return getParamVal(def.ArrayInfo, value)
		}
		val := make([]any, 0, len(array)+1)
		for _, value := range array {
			switch value.(type) {
			case map[string]Param:
				valMap := make(map[string]any, len(array)+1)
				for _, tp := range value.(map[string]Param) {
					valMap[tp.Identifier], err = tp.ToVal()
					if err != nil {
						return nil, err
					}
				}
				val = append(val, valMap)
			default:
				val = append(val, value)
			}
		}
		return val, nil
	default:
		return value, nil
	}
}

func (tp *Param) ToVal() (any, error) {
	if tp == nil {
		return nil, errors.Parameter.AddMsgf("Param is nil")
	}
	return getParamVal(tp.Define, tp.Value)
}

// 从设备参数中获取值
func GetVal(d *schema.Define, val any) (any, error) {
	switch d.Type {
	case schema.DataTypeBool:
		return cast.ToBoolE(val)
	case schema.DataTypeInt:
		if num, err := cast.ToInt64E(val); err != nil {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			if validateDataRange && (num > cast.ToInt64(d.Max) || num < cast.ToInt64(d.Min)) {
				return nil, errors.OutRange.AddDetailf("value %v out of range:[%s,%s]", val, d.Max, d.Min)
			}
			step := cast.ToInt64(d.Step)
			if step != 0 {
				num = num / step * step
			}
			return num, nil
		}
	case schema.DataTypeFloat:
		if num, err := cast.ToFloat64E(val); err != nil {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			if validateDataRange && (num > cast.ToFloat64(d.Max) || num < cast.ToFloat64(d.Min)) {
				return nil, errors.OutRange.AddDetailf(
					"value %v out of range:[%s,%s]", val, d.Max, d.Min)
			}
			step := cast.ToFloat64(d.Step)
			if step != 0 && !math.IsNaN(step) && !math.IsInf(step, 0) {
				num = math.Floor(num/step) * step
			}
			return num, nil
		}
	case schema.DataTypeString:
		if str, ok := val.(string); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			if validateDataRange && (len(str) > cast.ToInt(d.Max)) {
				return nil, errors.OutRange.AddDetailf("value %v out of range:%s", val, d.Max)
			}
			return str, nil
		}
	case schema.DataTypeEnum: //枚举类型 报文中传递的是数字
		if num, err := cast.ToInt64E(val); err != nil {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			_, ok := d.Mapping[cast.ToString(num)]
			if !ok {
				return nil, errors.OutRange.AddDetailf("value %v not in enum", val)
			}
			return num, nil
		}
	case schema.DataTypeTimestamp:
		switch val.(type) {
		case json.Number:
			ret, err := val.(json.Number).Int64()
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			return ret, nil
		case string:
			ret, err := cast.ToInt64E(val)
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			return ret, nil
		}
	case schema.DataTypeStruct:
		if strut, ok := val.(map[string]any); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			getParam := make(map[string]Param, len(strut))
			for k, v := range strut {
				sv, ok := d.Spec[k]
				if ok == false {
					continue
				}
				tp := Param{
					Identifier: sv.Identifier,
					Name:       sv.Name,
				}
				err := tp.SetByDefine(&sv.DataType, v)
				if err == nil {
					getParam[k] = tp
				} else if !errors.Cmp(err, errors.NotFind) {
					return nil, errors.Fmt(err).AddDetail(sv.Identifier)
				}
			}
			return getParam, nil
		}
	case schema.DataTypeArray:
		if arr, ok := val.([]any); !ok { //如果是指定id的方式
			return GetVal(d.ArrayInfo, val)
		} else {
			if len(arr) == 0 {
				return d, errors.NotFind
			}
			if len(arr) > cast.ToInt(d.Max) {
				errors.Parameter.AddDetailf("数组的长度超过物模型定义的:%v,上传了:%v", d.Max, len(arr))
			}
			getParam := make([]any, 0, len(arr)+1)
			for _, v := range arr {
				param, err := GetVal(d.ArrayInfo, v)
				if err == nil {
					getParam = append(getParam, param)
				} else if !errors.Cmp(err, errors.NotFind) {
					return nil, errors.Fmt(err).AddDetail(fmt.Sprint(v))
				}
			}
			return getParam, nil
		}
	}
	return nil, errors.Parameter.AddDetail("need param")
}

func ToParamValues(tp map[string]Param) (map[string]any, error) {
	ret := make(map[string]any, len(tp))
	var err error
	for k, v := range tp {
		ret[k], err = v.ToVal()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
