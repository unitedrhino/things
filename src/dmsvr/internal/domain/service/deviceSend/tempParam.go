package deviceSend

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
)

const (
	//是否校验参数范围
	validateDataRange = false
)

type TempParam struct {
	ID       string `json:"id"`       //标识符
	Name     string `json:"name"`     //功能名称
	Desc     string `json:"gesc"`     //描述
	Mode     string `json:"mode"`     //读写乐行:rw(可读可写) r(只读)
	Required bool   `json:"required"` //是否必须
	Type     string `json:"type"`     //事件类型: 信息:info  告警alert  故障:fault
	Value    struct {
		Type   schema.DataType   `json:"type"`              //参数类型:bool int string struct float timestamp array enum
		Maping map[string]string `json:"mapping,omitempty"` //枚举及bool类型:bool enum
		Min    string            `json:"min,omitempty"`     //数值最小值:int string float
		Max    string            `json:"max,omitempty"`     //数值最大值:int string float
		Start  string            `json:"start,omitempty"`   //初始值:int float
		Step   string            `json:"step,omitempty"`    //步长:int float
		Unit   string            `json:"unit,omitempty"`    //单位:int float
		Value  any               `json:"Value"`
		/*
			读到的数据  如果是是数组则类型为[]interface{}  如果是结构体类型则为map[id]TempParam
				interface 为数据内容  					string为结构体的key value 为数据内容
		*/
	} `json:"Value"` //数据定义
}

func (t *TempParam) AddDefine(d *schema.Define, val any) (err error) {
	t.Value.Type = d.Type
	t.Value.Type = d.Type
	t.Value.Maping = make(map[string]string)
	for k, v := range d.Maping {
		t.Value.Maping[k] = v
	}
	t.Value.Maping = d.Maping
	t.Value.Min = d.Min
	t.Value.Max = d.Max
	t.Value.Start = d.Start
	t.Value.Step = d.Step
	t.Value.Unit = d.Unit
	t.Value.Value, err = GetVal(d, val)
	return err
}

func ToVal(tp map[string]TempParam) map[string]any {
	ret := make(map[string]any, len(tp))
	for k, v := range tp {
		ret[k] = v.ToVal()
	}
	return ret
}

func (t *TempParam) ToVal() any {
	if t == nil {
		panic("TempParam is nil")
	}

	switch t.Value.Type {
	case schema.STRUCT:
		v, ok := t.Value.Value.(map[string]TempParam)
		if ok == false {
			return nil
		}
		val := make(map[string]any, len(v)+1)
		for _, tp := range v {
			val[tp.ID] = tp.ToVal()
		}
		return val
	case schema.ARRAY:
		array, ok := t.Value.Value.([]any)
		if ok == false {
			return nil
		}
		val := make([]any, 0, len(array)+1)
		for _, value := range array {
			switch value.(type) {
			case map[string]TempParam:
				valMap := make(map[string]any, len(array)+1)
				for _, tp := range value.(map[string]TempParam) {
					valMap[tp.ID] = tp.ToVal()
				}
				val = append(val, valMap)
			default:
				val = append(val, value)
			}
		}
		return val
	default:
		return t.Value.Value
	}
}

func GetVal(d *schema.Define, val any) (any, error) {
	switch d.Type {
	case schema.BOOL:
		switch val.(type) {
		case bool:
			return val.(bool), nil
		case json.Number:
			num := val.(json.Number).String()
			if num == "0" {
				return false, nil
			} else {
				return true, nil
			}
		}
	case schema.INT:
		if num, ok := val.(json.Number); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			ret, err := num.Int64()
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			if validateDataRange && (ret > cast.ToInt64(d.Max) || ret < cast.ToInt64(d.Min)) {
				return nil, errors.OutRange.AddDetailf("value %v out of range:[%s,%s]", val, d.Max, d.Min)
			}
			return ret, nil
		}
	case schema.FLOAT:
		if num, ok := val.(json.Number); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			ret, err := num.Float64()
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			if validateDataRange && (ret > cast.ToFloat64(d.Max) || ret < cast.ToFloat64(d.Min)) {
				return nil, errors.OutRange.AddDetailf(
					"value %v out of range:[%s,%s]", val, d.Max, d.Min)
			}
			return ret, nil
		}
	case schema.STRING:
		if str, ok := val.(string); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			if validateDataRange && (len(str) > cast.ToInt(d.Max)) {
				return nil, errors.OutRange.AddDetailf("value %v out of range:%s", val, d.Max)
			}
			return str, nil
		}
	case schema.ENUM: //枚举类型 报文中传递的是数字
		if num, ok := val.(json.Number); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			ret, err := num.Int64()
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			_, ok := d.Maping[string(num)]
			if !ok {
				return nil, errors.OutRange.AddDetailf("value %v not in enum", val)
			}
			return ret, nil
		}
	case schema.TIMESTAMP:
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
	case schema.STRUCT:
		if strut, ok := val.(map[string]any); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			getParam := make(map[string]TempParam, len(strut))
			for k, v := range strut {
				sv, ok := d.Spec[k]
				if ok == false {
					continue
				}
				tp := TempParam{
					ID:   sv.ID,
					Name: sv.Name,
				}
				err := tp.AddDefine(&sv.DataType, v)
				if err == nil {
					getParam[k] = tp
				} else if !errors.Cmp(err, errors.NotFind) {
					return nil, errors.Fmt(err).AddDetail(sv.ID)
				}
			}
			return getParam, nil
		}
	case schema.ARRAY:
		if arr, ok := val.([]any); !ok {
			return nil, errors.Parameter.AddDetail(fmt.Sprint(val))
		} else {
			if len(arr) == 0 {
				return d, errors.NotFind
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
