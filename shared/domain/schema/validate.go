package schema

import (
	"encoding/json"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
)

const (
	IDFomatChect     = false //是否检查ID是否以数字开头
	IDLen            = 128   //标识符的长度
	NameLen          = 128   //参数名称的长度
	DescLen          = 80    //描述的最大长度
	DefineMappingLen = 20
	DefineUnitLen    = 12
	DefineIntMax     = 9999999999999
	DefineIntMin     = -9999999999999
	DefineStringMax  = 2048
	DefineSpecsLen   = 10
	ParamsLen        = 20
)

func ValidateWithFmt(schemaStr []byte) (*Model, error) {
	schema := Model{}
	err := json.Unmarshal(schemaStr, &schema)
	if err != nil {
		return nil, errors.Parameter.WithMsg("模板的json格式不对")
	}
	err = schema.ValidateWithFmt()
	if err != nil {
		return nil, err
	}
	return schema.init(), err
}

func (t *Model) ValidateWithFmt() error {
	idMap := make(map[string]struct{}, len(t.Actions)+len(t.Events)+len(t.Properties))
	for i := range t.Properties {
		if _, ok := idMap[t.Properties[i].ID]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("属性的id重复了:%v", t.Properties[i].ID)
		}
		idMap[t.Properties[i].ID] = struct{}{}
		err := t.Properties[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	for i := range t.Events {
		if _, ok := idMap[t.Events[i].ID]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("属性的id重复了:%v", t.Events[i].ID)
		}
		idMap[t.Events[i].ID] = struct{}{}
		err := t.Events[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	for i := range t.Actions {
		if _, ok := idMap[t.Actions[i].ID]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("属性的id重复了:%v", t.Actions[i].ID)
		}
		idMap[t.Actions[i].ID] = struct{}{}
		err := t.Actions[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Action) ValidateWithFmt() error {
	if err := IDValidate(a.ID); err != nil {
		return err
	}
	if err := NameValidate(a.Name); err != nil {
		return err
	}
	if err := DescValidate(a.Desc); err != nil {
		return err
	}
	if err := a.Input.ValidateWithFmt(); err != nil {
		return err
	}
	return a.Output.ValidateWithFmt()
}

func (e *Event) ValidateWithFmt() error {
	if err := IDValidate(e.ID); err != nil {
		return err
	}
	if err := NameValidate(e.Name); err != nil {
		return err
	}
	if err := DescValidate(e.Desc); err != nil {
		return err
	}
	if e.Type != EventTypeInfo && e.Type != EventTypeAlert && e.Type != EventTypeFault {
		return errors.Parameter.WithMsgf("事件类型类型只能为info,alert及fault,收到:%v", e.Type)
	}
	return e.Params.ValidateWithFmt()
}

func (p *Property) ValidateWithFmt() error {
	if err := IDValidate(p.ID); err != nil {
		return err
	}
	if err := NameValidate(p.Name); err != nil {
		return err
	}
	if err := DescValidate(p.Desc); err != nil {
		return err
	}
	if p.Mode != PropertyModeWR && p.Mode != PropertyModeR {
		return errors.Parameter.WithMsgf("属性读写类型只能为wr及r,收到:%v", p.Mode)
	}
	return p.Define.ValidateWithFmt()
}

func IDValidate(id string) error {
	if len(id) > IDLen || len(id) == 0 {
		return errors.Parameter.WithMsgf(
			"标识符的第一个字符不能是数字，支持英文、数字、下划线的组合，最多不超过%v个字符,标识符:%v", IDLen, id)
	}
	if IDFomatChect {
		if !(id[0] <= '9' || id[0] >= '0') {
			return errors.Parameter.WithMsgf(
				"标识符的第一个字符不能是数字，支持英文、数字、下划线的组合，最多不超过%v个字符,标识符:%v", IDLen, id)
		}
	}
	return nil
}

func NameValidate(name string) error {
	if len([]rune(name)) > NameLen {
		return errors.Parameter.WithMsgf("名称支持中文、英文、数字、下划线的组合，最多不超过%v个字符,名称:%v", NameLen, name)
	}
	return nil
}
func DescValidate(desc string) error {
	if len([]rune(desc)) > DescLen {
		return errors.Parameter.WithMsgf("描述最多不超过%v个字符,描述:%v", DescLen, desc)
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
	return d.Specs.ValidateWithFmt()
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
	if err := NameValidate(s.Name); err != nil {
		return err
	}
	return s.DataType.ValidateWithFmt()
}
func (s Specs) ValidateWithFmt() error {
	if len(s) > DefineSpecsLen {
		return errors.Parameter.WithMsgf("结构体的参数最多只支持%v个", DefineSpecsLen)
	}
	specMap := make(map[string]struct{}, len(s))
	for i := range s {
		if _, ok := specMap[s[i].ID]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("结构体类型中的id重复了:%v", s[i].ID)
		}
		specMap[s[i].ID] = struct{}{}
		err := s[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Param) ValidateWithFmt() error {
	if err := IDValidate(p.ID); err != nil {
		return err
	}
	if err := NameValidate(p.Name); err != nil {
		return err
	}
	return p.Define.ValidateWithFmt()
}
func (p Params) ValidateWithFmt() error {
	if len(p) > ParamsLen {
		return errors.Parameter.WithMsgf("参数最多只支持%v个", ParamsLen)
	}
	paramMap := make(map[string]struct{}, len(p))
	for i := range p {
		if _, ok := paramMap[p[i].ID]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("参数的id重复了:%v", p[i].ID)
		}
		paramMap[p[i].ID] = struct{}{}
		err := p[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	return nil
}
