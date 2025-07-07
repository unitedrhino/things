package schema

import (
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"math"
)

func (d *Define) FmtValue(val any) (any, error) {
	switch d.Type {
	case DataTypeBool:
		return cast.ToInt64E(val)
	case DataTypeInt:
		if num, err := cast.ToInt64E(val); err != nil {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			step := cast.ToInt64(d.Step)
			if step != 0 {
				num = num / step * step
			}
			return num, nil
		}
	case DataTypeEnum, DataTypeTimestamp:
		if num, err := cast.ToInt64E(val); err != nil {
			num2, er := cast.ToFloat64E(val)
			if er != nil {
				return nil, errors.Parameter.AddDetail(err)
			}
			num, err = cast.ToInt64E(num2)
			if err != nil {
				return nil, errors.Parameter.AddDetail(err)
			}
			return num, nil
		} else {
			return num, nil
		}
	case DataTypeFloat:
		if num, err := cast.ToFloat64E(val); err != nil {
			return nil, errors.Parameter.AddDetail(val)
		} else {
			step := cast.ToFloat64(d.Step)
			if step != 0 && !math.IsNaN(step) && !math.IsInf(step, 0) {
				num = utils.StepFloat(num, step)
			}
			return num, nil
		}
	case DataTypeString:
		if str, err := cast.ToStringE(val); err != nil {
			return nil, errors.Parameter.AddDetail(err)
		} else {
			return str, nil
		}
	case DataTypeStruct:
		switch val.(type) {
		case map[string]any:
			strut := val.(map[string]any)
			var ret = map[string]any{}
			for k, v := range strut {
				sv, ok := d.Spec[k]
				if ok == false {
					continue
				}
				va, err := sv.DataType.FmtValue(v)
				if err != nil {
					return nil, err
				}
				ret[k] = va
			}
			return ret, nil
		case string, []byte: //需要json转换一下
			var ret = map[string]any{}
			var in []byte
			if v, ok := val.(string); ok {
				in = []byte(v)
			} else {
				in = val.([]byte)
			}
			err := utils.Unmarshal(in, &ret)
			if err != nil {
				return nil, errors.Parameter.AddDetail(err)
			}
			return d.FmtValue(ret)
		}
	case DataTypeArray:
		switch val.(type) {
		case []any:
			arr := val.([]any)
			if len(arr) == 0 {
				return d, errors.NotFind
			}
			getParam := make([]any, 0, len(arr)+1)
			for _, v := range arr {
				param, err := d.ArrayInfo.FmtValue(v)
				if err == nil {
					getParam = append(getParam, param)
				} else if !errors.Cmp(err, errors.NotFind) {
					return nil, errors.Fmt(err).AddDetail(fmt.Sprint(v))
				}
			}
			return getParam, nil
		case string, []byte:
			var ret []any
			var in []byte
			if v, ok := val.(string); ok {
				in = []byte(v)
			} else {
				in = val.([]byte)
			}
			err := utils.Unmarshal(in, &ret)
			if err != nil {
				return nil, errors.Parameter.AddDetail(err)
			}
			return d.FmtValue(ret)
		}
	}
	return nil, errors.Parameter.AddDetail("need param")
}
