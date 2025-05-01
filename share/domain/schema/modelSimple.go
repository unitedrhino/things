package schema

type (
	// Model 物模型协议-数据模板定义
	ModelSimple struct {
		Properties PropertiesSimple `json:"properties,omitempty"` //属性
		Events     EventsSimple     `json:"events,omitempty"`     //事件
		Actions    ActionsSimple    `json:"actions,omitempty"`    //行为
	}
	/*事件*/
	EventSimple struct {
		Identifier string       `json:"id"`     //标识符 (统一)
		Name       string       `json:"name"`   //功能名称
		Type       EventType    `json:"type"`   //事件类型: 1:信息:info  2:告警alert  3:故障:fault
		Params     ParamSimples `json:"params"` //事件参数
	}
	EventsSimple []EventSimple

	ParamSimple struct {
		Identifier string `json:"id"`   //参数标识符
		Name       string `json:"name"` //参数名称
		Define            //参数定义
	}
	ParamSimples []ParamSimple
	/*行为*/
	ActionSimple struct {
		Identifier string       `json:"id"`     //标识符 (统一)
		Name       string       `json:"name"`   //功能名称
		Dir        ActionDir    `json:"dir"`    //调用方向
		Input      ParamSimples `json:"input"`  //调用参数
		Output     ParamSimples `json:"output"` //返回参数
	}
	ActionsSimple []ActionSimple

	/*属性*/
	PropertySimple struct {
		Identifier string       `json:"id"`   //标识符 (统一)
		Name       string       `json:"name"` //功能名称
		Mode       PropertyMode `json:"mode"` //读写类型:rw(可读可写) r(只读)
		Define                  //数据定义
	}
	PropertiesSimple []PropertySimple
)

func (m *ModelSimple) ToModel() *Model {
	var ret Model
	for _, p := range m.Properties {
		ret.Properties = append(ret.Properties, Property{
			CommonParam: CommonParam{Identifier: p.Identifier, Name: p.Name},
			Mode:        p.Mode,
			Define:      p.Define,
		})
	}
	for _, e := range m.Events {
		ret.Events = append(ret.Events, Event{
			CommonParam: CommonParam{Identifier: e.Identifier, Name: e.Name},
			Type:        e.Type,
			Params:      e.Params.ToModel(),
		})
	}
	for _, a := range m.Actions {
		ret.Actions = append(ret.Actions, Action{
			CommonParam: CommonParam{Identifier: a.Identifier, Name: a.Name},
			Output:      a.Output.ToModel(),
			Input:       a.Input.ToModel(),
		})
	}
	return &ret
}
func (p ParamSimple) ToModel() Param {
	return Param{
		Identifier: p.Identifier,
		Name:       p.Name,
		Define:     p.Define,
	}
}
func (p ParamSimples) ToModel() Params {
	var ret Params
	for _, v := range p {
		ret = append(ret, v.ToModel())
	}
	return ret
}

// 转换成人看得懂的描述
func (m *ModelSimple) ToHumanDesc() string {
	return ""
}
