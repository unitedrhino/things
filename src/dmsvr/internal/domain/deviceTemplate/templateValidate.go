package deviceTemplate

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
)

const (
	IDLen            = 32 //标识符的长度
	NameLen          = 20 //参数名称的长度
	DefineMappingLen = 20
	DefineUnitLen    = 12
	DefineIntMax     = 9999999999999
	DefineIntMin     = -9999999999999
	DefineStringMax  = 2048
)

func IDValidate(id string) error {
	if len(id) > 32 || len(id) == 0 {
		return errors.Parameter.WithMsgf("标识符的第一个字符不能是数字，支持英文、数字、下划线的组合，最多不超过32个字符")
	}
	if id[0] <= '9' || id[0] >= '0' {
		return errors.Parameter.WithMsgf("标识符的第一个字符不能是数字，支持英文、数字、下划线的组合，最多不超过32个字符")
	}
	return nil
}

func (d *Define) ValidateWithFmt() error {
	switch d.Type {
	case BOOL:
		return d.ValidateWithFmtBool()
	case INT:
		return d.ValidateWithFmtInt()
	case STRING:
		return d.ValidateWithFmtString()
	case STRUCT:
		return d.ValidateWithFmtStruct()
	case FLOAT:
		return d.ValidateWithFmtFloat()
	case TIMESTAMP:
		return d.ValidateWithFmtTimeStamp()
	case ARRAY:
		return d.ValidateWithFmtArray()
	case ENUM:
		return d.ValidateWithFmtEnum()

	}
	return nil
}
func (d *Define) ValidateWithFmtBool() error {
	if len(d.Maping) != 2 {
		return errors.Parameter.WithMsgf("布尔的数据定义不正确:%v", d.Maping)
	}
	if v, ok := d.Maping["0"]; !ok {
		return errors.Parameter.WithMsgf("布尔的数据定义不正确:%v", d.Maping)
	} else {
		if len(v) > DefineMappingLen {
			return errors.Parameter.WithMsgf("布尔的0数据定义值长度过大:%v", d.Maping)
		}
	}
	if v, ok := d.Maping["1"]; !ok {
		return errors.Parameter.WithMsgf("布尔的数据定义不正确:%v", d.Maping)
	} else {
		if len(v) > DefineMappingLen {
			return errors.Parameter.WithMsgf("布尔的1数据定义值长度过大:%v", d.Maping)
		}
	}
	d.Min = ""
	d.Max = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil

	return nil
}
func (d *Define) ValidateWithFmtInt() error {
	max, err := cast.ToInt64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("整数的最大值定义不是数字:%v", d.Max)
	}
	if max > DefineIntMax {
		max = DefineIntMax
		d.Max = cast.ToString(max)
	}
	min, err := cast.ToInt64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("整数的最小值定义不是数字:%v", d.Min)
	}
	if min < DefineIntMin {
		min = DefineIntMin
		d.Min = cast.ToString(min)
	}
	if len(d.Unit) > DefineUnitLen {
		return errors.Parameter.WithMsgf("整数的单位定义值长度过大:%v", d.Unit)
	}
	step, err := cast.ToInt64E(d.Step)
	if err != nil {
		return errors.Parameter.WithMsgf("整数的单位定义不是数字:%v", d.Max)
	}
	if step > max {
		step = max
	}
	if step < 1 {
		step = 1
	}
	d.Maping = nil
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}
func (d *Define) ValidateWithFmtString() error {
	max, err := cast.ToInt64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("字符串的最大值定义不是数字:%v", d.Max)
	}
	if max > DefineStringMax {
		max = DefineStringMax
		d.Max = cast.ToString(max)
	}
	d.Min = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Maping = nil
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}
func (d *Define) ValidateWithFmtStruct() error {
	d.Max = ""
	d.Min = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Maping = nil
	d.ArrayInfo = nil
	for i := range d.Specs {
		err := d.Specs[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Define) ValidateWithFmtFloat() error {
	max, err := cast.ToFloat64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("浮点型的最大值定义不是数字:%v", d.Max)
	}
	if max > DefineIntMax {
		max = DefineIntMax
		d.Max = cast.ToString(max)
	}
	min, err := cast.ToFloat64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("浮点型的最小值定义不是数字:%v", d.Min)
	}
	if min < DefineIntMin {
		min = DefineIntMin
		d.Min = cast.ToString(min)
	}
	if len(d.Unit) > DefineUnitLen {
		return errors.Parameter.WithMsgf("浮点型的单位定义值长度过大:%v", d.Unit)
	}
	step, err := cast.ToFloat64E(d.Step)
	if err != nil {
		return errors.Parameter.WithMsgf("浮点型的单位定义不是数字:%v", d.Max)
	}
	if step > max {
		step = max
	}
	if step < 1 {
		step = 1
	}
	d.Maping = nil
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}
func (d *Define) ValidateWithFmtTimeStamp() error {
	d.Max = ""
	d.Min = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Maping = nil
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}
func (d *Define) ValidateWithFmtArray() error {
	d.Max = ""
	d.Min = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Maping = nil
	d.Specs = nil
	d.Spec = nil
	return d.ArrayInfo.ValidateWithFmt()
}
func (d *Define) ValidateWithFmtEnum() error {
	if len(d.Maping) == 0 {
		return errors.Parameter.WithMsgf("枚举的数据定义长度不能为0")
	}
	for k, v := range d.Maping {
		_, err := cast.ToInt64E(k)
		if err != nil {
			return errors.Parameter.WithMsgf("枚举的枚举键值定义不是数字:%v", k)
		}
		if len(v) > DefineMappingLen {
			return errors.Parameter.WithMsgf("枚举的%v数据定义值长度过大:%v", k, v)
		}
	}
	d.Min = ""
	d.Max = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}

func (s *Spec) ValidateWithFmt() error {
	if err := IDValidate(s.ID); err != nil {
		return err
	}
	if len(s.Name) > NameLen {
		return errors.Parameter.WithMsg("支持中文、英文、数字、下划线的组合，最多不超过20个字符")
	}
	return s.DataType.ValidateWithFmt()
}
