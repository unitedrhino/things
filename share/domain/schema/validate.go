package schema

import (
	"encoding/json"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/share/errors"
	"github.com/spf13/cast"
)

const (
	IDFFormatCheck   = true //是否检查ID是否以数字开头
	IDLen            = 128  //标识符的长度
	NameLen          = 128  //参数名称的长度
	DescLen          = 80   //描述的最大长度
	DefineMappingLen = 50
	DefineUnitLen    = 12
	DefineIntMax     = 9999999999999
	DefineIntMin     = -9999999999999
	DefineStringMax  = 2048
	DefineArrayMax   = 1024
	DefineSpecsLen   = 20
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

func (m *Model) Aggregation(m2 *Model) *Model {
	if m == nil {
		return m2
	}
	if m2 == nil {
		return m
	}
	for _, v := range m2.Properties {
		vv := v
		vv.Define.Spec = nil
		m.Properties = append(m.Properties, v)
	}
	for _, v := range m2.Events {
		m.Events = append(m.Events, v)
	}
	for _, v := range m2.Actions {
		m.Actions = append(m.Actions, v)
	}
	return m
}

func (m *Model) Copy() *Model {
	if m == nil {
		return nil
	}
	newOne := Model{Profile: m.Profile}
	for _, p := range m.Properties {
		v := p
		v.Define = v.Define.Copy()
		newOne.Properties = append(newOne.Properties, v)

	}
	for _, e := range m.Events {
		v := e
		v.Param = nil
		v.Params = nil
		for _, p := range e.Params {
			vv := p
			vv.Define = vv.Define.Copy()
			v.Params = append(v.Params, vv)
		}
		newOne.Events = append(newOne.Events, v)

	}
	for _, a := range m.Actions {
		v := a
		v.In = nil
		v.Out = nil
		v.Output = nil
		v.Input = nil
		for _, p := range a.Input {
			vv := p
			vv.Define = vv.Define.Copy()
			v.Input = append(v.Input, vv)
		}
		for _, p := range a.Output {
			vv := p
			vv.Define = vv.Define.Copy()
			v.Output = append(v.Output, vv)
		}
		newOne.Actions = append(newOne.Actions, v)
	}
	return &newOne
}

func (m *Model) ValidateWithFmt() error {
	idMap := make(map[string]struct{}, len(m.Actions)+len(m.Events)+len(m.Properties))
	for i := range m.Properties {
		if _, ok := idMap[m.Properties[i].Identifier]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("属性的id重复了:%v", m.Properties[i].Identifier)
		}
		idMap[m.Properties[i].Identifier] = struct{}{}
		err := m.Properties[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	for i := range m.Events {
		if _, ok := idMap[m.Events[i].Identifier]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("属性的id重复了:%v", m.Events[i].Identifier)
		}
		idMap[m.Events[i].Identifier] = struct{}{}
		err := m.Events[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	for i := range m.Actions {
		if _, ok := idMap[m.Actions[i].Identifier]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("属性的id重复了:%v", m.Actions[i].Identifier)
		}
		idMap[m.Actions[i].Identifier] = struct{}{}
		err := m.Actions[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	m.init()

	return nil
}

func (a *Action) ValidateWithFmt() error {
	if err := IDValidate(a.Identifier); err != nil {
		return err
	}
	if err := NameValidate(a.Name); err != nil {
		return err
	}
	if !utils.SliceIn(a.Dir, ActionDirDown, ActionDirUp, "") {
		return errors.Parameter.WithMsgf("行为的控制方向只能为up及down,收到:%v", a.Dir)
	}
	if a.Dir == "" {
		a.Dir = ActionDirDown
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
	if err := IDValidate(e.Identifier); err != nil {
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
	if err := IDValidate(p.Identifier); err != nil {
		return err
	}
	if err := NameValidate(p.Name); err != nil {
		return err
	}
	if err := DescValidate(p.Desc); err != nil {
		return err
	}
	if p.Mode == "" {
		p.Mode = PropertyModeRW
	}
	if p.IsPassword == 0 {
		p.IsPassword = def.False
	}
	if p.Mode != PropertyModeRW && p.Mode != PropertyModeR {
		return errors.Parameter.WithMsgf("属性读写类型只能为rw及r,收到:%v", p.Mode)
	}
	return p.Define.ValidateWithFmt()
}

func IDValidate(id string) error {
	if len(id) > IDLen || len(id) == 0 {
		return errors.Parameter.WithMsgf(
			"标识符的第一个字符不能是数字，支持英文、数字、下划线的组合，最多不超过%v个字符,标识符:%v", IDLen, id)
	}
	if IDFFormatCheck {
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
	case DataTypeBool:
		return d.ValidateWithFmtBool()
	case DataTypeInt:
		return d.ValidateWithFmtInt()
	case DataTypeString:
		return d.ValidateWithFmtString()
	case DataTypeStruct:
		return d.ValidateWithFmtStruct()
	case DataTypeFloat:
		return d.ValidateWithFmtFloat()
	case DataTypeTimestamp:
		return d.ValidateWithFmtTimeStamp()
	case DataTypeArray:
		return d.ValidateWithFmtArray()
	case DataTypeEnum:
		return d.ValidateWithFmtEnum()
	}
	return errors.Parameter.WithMsgf("定义的类型不支持:%v", d.Type)
}
func (d *Define) ValidateWithFmtBool() error {
	if len(d.Mapping) == 0 {
		d.Mapping = map[string]string{
			"0": "关",
			"1": "开",
		}
	}
	if len(d.Mapping) != 2 {
		return errors.Parameter.WithMsgf("布尔的数据定义不正确:%v", d.Mapping)
	}
	if v, ok := d.Mapping["0"]; !ok {
		return errors.Parameter.WithMsgf("布尔的数据定义不正确:%v", d.Mapping)
	} else {
		if len(v) > DefineMappingLen {
			return errors.Parameter.WithMsgf("布尔的0数据定义值长度过大:%v", d.Mapping)
		}
	}
	if v, ok := d.Mapping["1"]; !ok {
		return errors.Parameter.WithMsgf("布尔的数据定义不正确:%v", d.Mapping)
	} else {
		if len(v) > DefineMappingLen {
			return errors.Parameter.WithMsgf("布尔的1数据定义值长度过大:%v", d.Mapping)
		}
	}
	_, ok := d.Mapping[d.Start]
	if d.Start != "" && !ok {
		return errors.Parameter.WithMsgf("布尔的初始值没有在定义范围内:%v", d.Start)
	}
	d.Min = ""
	d.Max = ""
	d.Step = ""
	d.Unit = ""
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil

	return nil
}
func (d *Define) ValidateWithFmtInt() error {
	if d.Max == "" {
		d.Max = cast.ToString(DefineIntMax)
	}
	max, err := cast.ToInt64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("整数的最大值定义不是数字:%v", d.Max)
	}
	if max > DefineIntMax {
		max = DefineIntMax
		d.Max = cast.ToString(max)
	}
	if d.Min == "" {
		d.Min = cast.ToString(DefineIntMin)
	}
	min, err := cast.ToInt64E(d.Min)
	if err != nil {
		return errors.Parameter.WithMsgf("整数的最小值定义不是数字:%v", d.Min)
	}
	if min < DefineIntMin {
		min = DefineIntMin
		d.Min = cast.ToString(min)
	}
	if max < min {
		return errors.Parameter.WithMsgf("整数的最大值需要大于最小值")
	}
	if len(d.Unit) > DefineUnitLen {
		return errors.Parameter.WithMsgf("整数的单位定义值长度过大:%v", d.Unit)
	}
	if d.Step == "" {
		d.Step = "1"
	}
	step, err := cast.ToInt64E(d.Step)
	if err != nil {
		return errors.Parameter.WithMsgf("整数的步长定义值类型不是数字:%v", d.Step)
	}
	if step > max {
		d.Step = cast.ToString(max)
	}
	if step < 1 {
		d.Step = cast.ToString(1)
	}

	d.Mapping = nil
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}
func (d *Define) ValidateWithFmtString() error {
	if d.Max == "" {
		d.Max = cast.ToString(DefineStringMax)
	}
	max, err := cast.ToInt64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("字符串的最大值定义不是数字类型:%v", d.Max)
	}
	if max > DefineStringMax {
		max = DefineStringMax
		d.Max = cast.ToString(max)
	}
	if len(d.Start) > int(max) {
		return errors.Parameter.WithMsgf("字符串的默认值不能大于最大值:%v", d.Start)
	}
	d.Step = ""
	d.Min = ""
	d.Unit = ""
	d.Mapping = nil
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
	d.Mapping = nil
	d.ArrayInfo = nil
	return d.Specs.ValidateWithFmt()
}
func (d *Define) ValidateWithFmtFloat() error {
	if d.Max == "" {
		d.Max = cast.ToString(DefineIntMax)
	}
	max, err := cast.ToFloat64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("浮点型的最大值定义不是数字类型:%v", d.Max)
	}
	if max > DefineIntMax {
		max = DefineIntMax
		d.Max = cast.ToString(max)
	}
	if d.Min == "" {
		d.Min = cast.ToString(DefineIntMin)
	}
	min, err := cast.ToFloat64E(d.Min)
	if err != nil {
		return errors.Parameter.WithMsgf("浮点型的最小值定义不是数字类型:%v", d.Min)
	}
	if min < DefineIntMin {
		min = DefineIntMin
		d.Min = cast.ToString(min)
	}
	if max < min {
		return errors.Parameter.WithMsgf("浮点型的最大值需要大于最小值")
	}
	if len(d.Unit) > DefineUnitLen {
		return errors.Parameter.WithMsgf("浮点型的单位定义值长度过大:%v", d.Unit)
	}
	if d.Step == "" {
		d.Step = "0.001"
	}
	step, err := cast.ToFloat64E(d.Step)
	if err != nil {
		return errors.Parameter.WithMsgf("浮点型的步长定义不是数字类型:%v", d.Step)
	}
	if step > max {
		d.Step = cast.ToString(max)
	}
	if step <= 0 {
		d.Step = cast.ToString(1)
	}

	d.Mapping = nil
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
	d.Mapping = nil
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}
func (d *Define) ValidateWithFmtArray() error {
	max, err := cast.ToInt64E(d.Max)
	if err != nil {
		return errors.Parameter.WithMsgf("数组类型的个数定义不是数字类型:%v", d.Max)
	}
	if max > DefineArrayMax {
		max = DefineArrayMax
		d.Max = cast.ToString(max)
	}
	if d.Max == "0" {
		return errors.Parameter.WithMsg("数组类型的个数定义不能小于0")
	}
	d.Min = ""
	d.Start = ""
	d.Step = ""
	d.Unit = ""
	d.Mapping = nil
	d.Specs = nil
	d.Spec = nil
	if d.ArrayInfo == nil {
		return errors.Parameter.WithMsgf("数组类型缺失arrayInfo结构体")
	}
	return d.ArrayInfo.ValidateWithFmt()
}
func (d *Define) ValidateWithFmtEnum() error {
	if len(d.Mapping) == 0 {
		return errors.Parameter.WithMsgf("枚举的数据定义长度不能为0")
	}
	for k, v := range d.Mapping {
		_, err := cast.ToInt64E(k)
		if err != nil {
			return errors.Parameter.WithMsgf("枚举的枚举键值定义不是数字:%v", k)
		}
		if len(v) > DefineMappingLen {
			return errors.Parameter.WithMsgf("枚举的%v数据定义值长度过大:%v", k, v)
		}
	}
	_, ok := d.Mapping[d.Start]
	if d.Start != "" && !ok {
		return errors.Parameter.WithMsgf("枚举的初始值没有在定义范围内:%v", d.Start)
	}
	d.Min = ""
	d.Max = ""
	d.Step = ""
	d.Unit = ""
	d.Specs = nil
	d.ArrayInfo = nil
	d.Spec = nil
	return nil
}

func (s *Spec) ValidateWithFmt() error {
	if err := IDValidate(s.Identifier); err != nil {
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
		if _, ok := specMap[s[i].Identifier]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("结构体类型中的id重复了:%v", s[i].Identifier)
		}
		specMap[s[i].Identifier] = struct{}{}
		err := s[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Param) ValidateWithFmt() error {
	if err := IDValidate(p.Identifier); err != nil {
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
		if _, ok := paramMap[p[i].Identifier]; ok {
			//如果有重复的需要返回错误
			return errors.Parameter.WithMsgf("参数的id重复了:%v", p[i].Identifier)
		}
		paramMap[p[i].Identifier] = struct{}{}
		err := p[i].ValidateWithFmt()
		if err != nil {
			return err
		}
	}
	return nil
}
