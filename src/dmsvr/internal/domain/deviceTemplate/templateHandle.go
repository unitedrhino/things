package deviceTemplate

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
)

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

func (t *Template) VerifyReqParam(dq DeviceReq, tt TEMP_TYPE) (map[string]TempParam, error) {
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
			return nil, errors.Parameter.AddDetail("need right eventId")
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
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case ACTION_INPUT:
		p, ok := t.Action[dq.ActionID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.In {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := dq.Params[v.ID]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case ACTION_OUTPUT:
		p, ok := t.Action[dq.ActionID]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.In {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := dq.Params[v.ID]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	}
	return getParam, nil
}

func (t *Template) VerifyRespParam(dr DeviceResp, id string, tt TEMP_TYPE) (map[string]TempParam, error) {
	getParam := make(map[string]TempParam, len(dr.Response))
	switch tt {
	case ACTION_OUTPUT:
		p, ok := t.Action[id]
		if ok == false {
			return nil, errors.Parameter.AddDetail("need right ActionID")
		}
		for k, v := range p.Out {
			tp := TempParam{
				ID:   v.ID,
				Name: v.Name,
			}
			param, ok := dr.Response[v.ID]
			if ok == false {
				return nil, errors.Parameter.AddDetail("need param:" + k)
			}
			err := tp.AddDefine(&v.Define, param)
			if err == nil {
				getParam[k] = tp
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, errors.Fmt(err).AddDetail(p.ID)
			}
		}
	}
	return getParam, nil
}
