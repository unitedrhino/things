package schema

import "encoding/json"

type (
	// Model 物模型协议-数据模板定义
	Model struct {
		Version    string      `json:"version"`    //版本号
		Properties Properties  `json:"properties"` //属性
		Events     Events      `json:"events"`     //事件
		Actions    Actions     `json:"actions"`    //行为
		Profile    Profile     `json:"profile"`    //配置信息
		Property   PropertyMap `json:"-"`          //内部使用,使用map加速匹配,key为id
		Event      EventMap    `json:"-"`          //内部使用,使用map加速匹配,key为id
		Action     ActionMap   `json:"-"`          //内部使用,使用map加速匹配,key为id
	}

	/*配置信息*/
	Profile struct {
		ProductID string `json:"productID"` //产品ID
	}

	/*内部使用*/
	PropertyMap map[string]*Property
	EventMap    map[string]*Event
	ActionMap   map[string]*Action

	/*结构体说明*/
	Spec struct {
		Identifier string `json:"identifier"` //参数标识符
		Name       string `json:"name"`       //参数名称
		DataType   Define `json:"dataType"`   //参数定义
	}
	Specs []Spec

	/*参数*/
	Param struct {
		Identifier string `json:"identifier"`       //参数标识符
		Name       string `json:"name"`             //参数名称
		Define     Define `json:"define,omitempty"` //参数定义
	}
	Params []Param

	/*事件*/
	Event struct {
		Identifier string            `json:"identifier"` //标识符
		Name       string            `json:"name"`       //功能名称
		Desc       string            `json:"desc"`       //描述
		Type       EventType         `json:"type"`       //事件类型: 1:信息:info  2:告警alert  3:故障:fault
		Params     Params            `json:"params"`     //事件参数
		Required   bool              `json:"required"`   //是否必须
		Param      map[string]*Param `json:"-"`          //内部使用,使用map加速匹配,key为id
	}
	Events []Event

	/*行为*/
	Action struct {
		Identifier string            `json:"identifier"` //标识符
		Name       string            `json:"name"`       //功能名称
		Desc       string            `json:"desc"`       //描述
		Input      Params            `json:"input"`      //调用参数
		Output     Params            `json:"output"`     //返回参数
		Required   bool              `json:"required"`   //是否必须
		In         map[string]*Param `json:"-"`          //内部使用,使用map加速匹配,key为id
		Out        map[string]*Param `json:"-"`          //内部使用,使用map加速匹配,key为id
	}
	Actions []Action

	/*属性*/
	Property struct {
		Identifier string       `json:"identifier"` //标识符
		Name       string       `json:"name"`       //功能名称
		Desc       string       `json:"desc"`       //描述
		Mode       PropertyMode `json:"mode"`       //读写类型:rw(可读可写) r(只读)
		Define     Define       `json:"define"`     //数据定义
		Required   bool         `json:"required"`   //是否必须
	}
	Properties []Property

	/*数据类型定义*/
	Define struct {
		Type      DataType          `json:"type"`                //参数类型:bool int string struct float timestamp array enum
		Maping    map[string]string `json:"mapping,omitempty"`   //枚举及bool类型:bool enum
		Min       string            `json:"min,omitempty"`       //数值最小值:int  float
		Max       string            `json:"max,omitempty"`       //数值最大值:int string float
		Start     string            `json:"start,omitempty"`     //初始值:int float
		Step      string            `json:"step,omitempty"`      //步长:int float
		Unit      string            `json:"unit,omitempty"`      //单位:int float
		Specs     Specs             `json:"specs,omitempty"`     //结构体:struct
		ArrayInfo *Define           `json:"arrayInfo,omitempty"` //数组:array
		Spec      map[string]*Spec  `json:"-"`                   //内部使用,使用map加速匹配,key为id
	}
)

func (m *Model) String() string {
	tls, _ := json.Marshal(m)
	return string(tls)
}

func (p *PropertyMap) GetIDs() []string {
	var ids []string
	for _, v := range *p {
		ids = append(ids, v.Identifier)
	}
	return ids
}
