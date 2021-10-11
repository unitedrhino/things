package device

import (
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"github.com/spf13/cast"
)

type TempParam struct {
	ID       string `json:"id"`       //标识符
	Name     string `json:"name"`     //功能名称
	Desc     string `json:"gesc"`     //描述
	Mode     string `json:"mode"`     //读写乐行:rw(可读可写) r(只读)
	Required bool   `json:"required"` //是否必须
	Type     string `json:"type"`     //事件类型: 信息:info  告警alert  故障:fault
	Value    struct {
		Type   string            `json:"type"`              //参数类型:bool int string struct float timestamp array enum
		Maping map[string]string `json:"mapping,omitempty"` //枚举及bool类型:bool enum
		Min    string            `json:"min,omitempty"`     //数值最小值:int string float
		Max    string            `json:"max,omitempty"`     //数值最大值:int string float
		Start  string            `json:"start,omitempty"`   //初始值:int float
		Step   string            `json:"step,omitempty"`    //步长:int float
		Unit   string            `json:"unit,omitempty"`    //单位:int float
		Value  interface{}       `json:"Value"`
		/*
			读到的数据  如果是是数组则类型为[]interface{}  如果是结构体类型则为map[id]TempParam
				interface 为数据内容  					string为结构体的key value 为数据内容
		*/
	} `json:"Value"` //数据定义
}

func (t *TempParam) AddDefine(d *Define, val interface{}) (err error) {
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
	//todo
	t.Value.Value, err = d.GetVal(val)
	return err
}

func ToVal(tp map[string]TempParam) map[string]interface{} {
	ret := make(map[string]interface{}, len(tp))
	for k, v := range tp {
		ret[k] = v.ToVal()
	}
	return ret
}

func (t *TempParam) ToVal() interface{} {
	if t == nil {
		panic("TempParam is nil")
	}

	switch t.Value.Type {
	case STRUCT:
		v, ok := t.Value.Value.(map[string]TempParam)
		if ok == false {
			return nil
		}
		val := make(map[string]interface{}, len(v)+1)
		for _, tp := range v {
			val[tp.ID] = tp.ToVal()
		}
		return val
	case ARRAY:
		array, ok := t.Value.Value.([]interface{})
		if ok == false {
			return nil
		}
		val := make([]interface{}, 0, len(array)+1)
		for _, value := range array {
			switch value.(type) {
			case map[string]TempParam:
				valMap := make(map[string]interface{}, len(array)+1)
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

func (d *Define) GetVal(val interface{}) (interface{}, error) {
	switch d.Type {
	case BOOL:
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
	case INT:
		if num, ok := val.(json.Number); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			ret, err := num.Int64()
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			if ret > cast.ToInt64(d.Max) || ret < cast.ToInt64(d.Min) {
				return nil, errors.OutRange.AddDetailf("value %v out of range:[%s,%s]", val, d.Max, d.Min)
			}
			return ret, nil
		}
	case FLOAT:
		if num, ok := val.(json.Number); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			ret, err := num.Float64()
			if err != nil {
				return nil, errors.Parameter.AddDetail(val)
			}
			if ret > cast.ToFloat64(d.Max) || ret < cast.ToFloat64(d.Min) {
				return nil, errors.OutRange.AddDetailf(
					"value %v out of range:[%s,%s]", val, d.Max, d.Min)
			}
			return ret, nil
		}
	case STRING:
		if str, ok := val.(string); !ok {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			if len(str) > cast.ToInt(d.Max) {
				return nil, errors.OutRange.AddDetailf("value %v out of range:%s", val, d.Max)
			}
			return str, nil
		}
	case ENUM: //枚举类型 报文中传递的是数字
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
	case TIMESTAMP:
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
	case STRUCT:
		if strut, ok := val.(map[string]interface{}); !ok {
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
	case ARRAY:
		if arr, ok := val.([]interface{}); !ok {
			return nil, errors.Parameter.AddDetail(fmt.Sprint(val))
		} else {
			if len(arr) == 0 {
				return d, errors.NotFind
			}
			getParam := make([]interface{}, 0, len(arr)+1)
			for _, v := range arr {
				param, err := d.ArrayInfo.GetVal(v)
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

func (t *Template) VerifyParam(dq DeviceReq, tt TEMP_TYPE) (map[string]TempParam, error) {
	if len(dq.Params) == 0 {
		return nil, errors.Parameter.AddDetail("need add params")
	}
	getParam := make(map[string]TempParam, len(dq.Params))
	switch tt {
	case PROPERTY:
		for k, v := range dq.Params {
			p, ok := t.Property[k]
			if ok == false {
				continue
			}
			tp := TempParam{
				ID:       p.ID,
				Name:     p.Name,
				Desc:     p.Desc,
				Mode:     p.Mode,
				Required: p.Required,
			}
			err := tp.AddDefine(&p.Define, v)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case EVENT:
		p, ok := t.Event[dq.EventID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need add eventId")
		}
		if dq.Type != p.Type {
			return nil, errors.Parameter.AddDetail("err type:" + dq.Type)
		}

		for k, v := range p.Param {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := dq.Params[k]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case ACTION_INPUT:
	case ACTION_OUTPUT:

	}
	return getParam, nil
}
