package dict

import (
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"github.com/spf13/cast"
)

func (d *Define) AddVal(val interface{}) (*Define ,error){
	switch d.Type {
	case BOOL:
		switch val.(type) {
		case bool:
			d.Val = val.(bool)
			return d,nil
		case json.Number:
			num :=val.(json.Number).String()
			if  num == "0" {
				d.Val = false
				return d,nil
			}else {
				d.Val = true
				return d,nil
			}
		}
	case INT:
		if num,ok:=val.(json.Number);!ok{
			return d, errors.Parameter.AddDetail(val)
		}else {
			ret,err := num.Int64()
			if err != nil {
				return d, errors.Parameter.AddDetail(val)
			}
			d.Val = ret
			return d,nil
		}
	case FLOAT:
		if num,ok:=val.(json.Number);!ok{
			return d, errors.Parameter.AddDetail(val)
		}else {
			ret,err := num.Float64()
			if err != nil {
				return d, errors.Parameter.AddDetail(val)
			}
			d.Val = ret
			return d,nil
		}
	case STRING:
		if str,ok:=val.(string);!ok{
			return d, errors.Parameter.AddDetail(val)
		}else {
			d.Val = str
			return d,nil
		}
	case ENUM: //枚举类型 报文中传递的是数字
		if num,ok:=val.(json.Number);!ok{
			return d, errors.Parameter.AddDetail(val)
		}else {
			ret,err := num.Int64()
			if err != nil {
				return d, errors.Parameter.AddDetail(val)
			}
			d.Val = ret
			return d,nil
		}
	case TIMESTAMP:
		switch val.(type) {
		case json.Number:
			ret,err := val.(json.Number).Int64()
			if err != nil {
				return d, errors.Parameter.AddDetail(val)
			}
			d.Val = ret
			return d,nil
		case string:
			ret,err := cast.ToInt64E(val)
			if err != nil {
				return d, errors.Parameter.AddDetail(val)
			}
			d.Val = ret
			return d,nil
		}
	case STRUCT:
		if strut,ok := val.(map[string]interface {});!ok{
			return d, errors.Parameter.AddDetail(val)
		}else {
			getParam := make(map[string]*Define,len(strut))
			for k,v :=range strut {
				sv,ok := d.Spec[k]
				if ok == false {
					continue
				}
				param,err := sv.DataType.AddVal(v)
				if err == nil {
					getParam[k] = param
				}else if !errors.Cmp(err,errors.NotFind) {
					return d,errors.Fmt(err).AddDetail(sv.ID)
				}
			}
			d.Val = getParam
			return d,nil
		}
	case ARRAY:
		if arr,ok := val.([]interface {});!ok{
			return d, errors.Parameter.AddDetail(fmt.Sprint(val))
		}else {
			if len(arr) == 0{
				return d, errors.NotFind
			}
			getParam := make([]*Define,0,len(arr)+1)
			for _,v :=range arr {
				param,err := d.ArrayInfo.AddVal(v)
				if err == nil {
					getParam = append(getParam,param)
				}else if !errors.Cmp(err,errors.NotFind) {
					return d,errors.Fmt(err).AddDetail(fmt.Sprint(v))
				}
			}
			d.Val = getParam
			return d, nil
		}
	}
	return d, errors.Parameter.AddDetail("need param")
}

func (t *Template)VerifyParam(param map[string]interface{},tt TEMP_TYPE)(map[string]interface{} ,error){
	if len(param) == 0 {
		return nil,errors.Parameter.AddDetail("need add params")
	}
	getParam := make(map[string]interface{},len(param))
	switch tt {
	case PROPERTY:
		for k,v := range param{
			p,ok := t.Property[k]
			if ok == false {
				continue
			}
			def := p.Define
			data,err := def.AddVal(v)
			if err == nil {
				getParam[p.ID] = data
			}else if !errors.Cmp(err,errors.NotFind) {
				return nil,errors.Fmt(err).AddDetail(p.ID)
			}
		}
	case ACTION_INPUT:
	case ACTION_OUTPUT:
	case EVENT:

	}
	return getParam,nil
}

