package dict

import "encoding/json"

//数据类型
const (
	BOOL = "bool"
	INT = "int"
	STRING = "string"
	STRUCT = "struct"
	FLOAT = "float"
	TIMESTAMP = "timestamp"
	ARRAY = "array"
	ENUM = "enum"
)

type TEMP_TYPE int64

const (
	PROPERTY TEMP_TYPE = iota
	ACTION_INPUT
	ACTION_OUTPUT
	EVENT
)

type Profile struct {
	ProductID  string `json:"ProductId"`  //产品ID
	CategoryID string `json:"CategoryId"` //产品种类ID
}


/*结构体说明*/
type Spec struct {
	ID       string `json:"id"`       //参数标识符
	Name     string `json:"name"`     //参数名称
	DataType Define `json:"dataType"` //参数定义
}

/*数据类型定义*/
type Define struct {
	Type      string            `json:"type"`                //参数类型:bool int string struct float timestamp array enum
	Maping    map[string]string `json:"mapping,omitempty"`   //枚举及bool类型:bool enum
	Min       string            `json:"min,omitempty"`       //数值最小值:int string float
	Max       string            `json:"max,omitempty"`       //数值最大值:int string float
	Start     string            `json:"start,omitempty"`     //初始值:int float
	Step      string            `json:"step,omitempty"`      //步长:int float
	Unit      string            `json:"unit,omitempty"`      //单位:int float
	Specs     []Spec            `json:"specs,omitempty"`     //结构体:struct
	ArrayInfo *Define           `json:"arrayInfo,omitempty"` //数组:array
	Spec     map[string]*Spec  	//内部使用,使用map加速匹配,key为id
	/*
	读到的数据  如果是是数组则类型为[]interface{}  如果是结构体类型则为map[string]interface{}
		interface 为数据内容  					string为结构体的key value 为数据内容
	*/
	Val      interface{}
}
func (d *Define)init()*Define{
	if d.Specs!= nil {
		d.Spec = make(map[string]*Spec,len(d.Specs)+1)
		for i:=0;i < len(d.Specs);i++{
			p := &d.Specs[i]
			d.Spec[p.ID] = p
		}
	}
	if d.ArrayInfo != nil {
		d.ArrayInfo.init()
	}
	return d
}

/*参数*/
type Param struct {
	ID     string `json:"id"`     //参数标识符
	Name   string `json:"name"`   //参数名称
	Define Define `json:"define,omitempty"` //参数定义
}

func (p *Param)init()*Param{
	p.Define.init()
	return p
}
/*事件*/
type Event struct {
	ID       string  `json:"id"`       //标识符
	Name     string  `json:"name"`     //功能名称
	Desc     string  `json:"desc"`     //描述
	Type     string  `json:"type"`     //事件类型: 信息:info  告警alert  故障:fault
	Params   []Param `json:"params"`   //事件参数
	Required bool    `json:"required"` //是否必须
	Param     map[string]*Param  `json:"-"` 		//内部使用,使用map加速匹配,key为id
}
func (e *Event)init()*Event{
	e.Param = make(map[string]*Param,len(e.Params)+1)
	for i:=0;i < len(e.Params);i++{
		p := &e.Params[i]
		p.init()
		e.Param[p.ID] = p
	}
	return e
}

/*行为*/
type Action struct {
	ID       string  `json:"id"`       //标识符
	Name     string  `json:"name"`     //功能名称
	Desc     string  `json:"desc"`     //描述
	Input    []Param `json:"input"`    //调用参数
	Output   []Param `json:"output"`   //返回参数
	Required bool    `json:"required"` //是否必须
	In     map[string]*Param  `json:"-"` 		//内部使用,使用map加速匹配,key为id
	Out     map[string]*Param  `json:"-"` 		//内部使用,使用map加速匹配,key为id
}
func (a *Action)init()*Action{
	a.In = make(map[string]*Param,len(a.Input)+1)
	a.Out = make(map[string]*Param,len(a.Output)+1)
	for i:=0;i < len(a.Input);i++{
		p := &a.Input[i]
		a.In[p.ID] = p
		p.init()
	}
	for i:=0;i < len(a.Output);i++{
		p := &a.Output[i]
		a.Out[p.ID] = p
		p.init()
	}
	return a
}

/*属性*/
type Property struct {
	ID       string `json:"id"`       //标识符
	Name     string `json:"name"`     //功能名称
	Desc     string `json:"gesc"`     //描述
	Mode     string `json:"mode"`     //读写乐行:rw(可读可写) r(只读)
	Define   Define `json:"define"`   //数据定义
	Required bool   `json:"required"` //是否必须
}
func (p *Property)init()*Property{
	p.Define.init()
	return p
}
/*数据模板定义*/
type Template struct {
	Version    string             `json:"version"`    //版本号
	Properties []Property          `json:"properties"` //属性
	Events     []Event            `json:"events"`     //事件
	Actions    []Action           `json:"actions"`    //行为
	Profile    Profile            `json:"profile"`    //配置信息
	Property   map[string]*Property `json:"-"` //内部使用,使用map加速匹配,key为id
	Event      map[string]*Event   `json:"-"` //内部使用,使用map加速匹配,key为id
	Action     map[string]*Action  `json:"-"` //内部使用,使用map加速匹配,key为id
}

func (t *Template)init()*Template{
	t.Property = make(map[string]*Property,len(t.Properties)+1)
	t.Event = make(map[string]*Event,len(t.Events)+1)
	t.Action = make(map[string]*Action,len(t.Actions)+1)
	for i:=0;i < len(t.Properties);i++{
		p := &t.Properties[i]
		t.Property[p.ID] = p
		p.init()
	}
	for i:=0;i < len(t.Events);i++{
		p := &t.Events[i]
		t.Event[p.ID] = p
		p.init()
	}
	for i:=0;i < len(t.Actions);i++{
		p := &t.Actions[i]
		t.Action[p.ID] = p
		p.init()
	}
	return t
}

func NewTemplate(templateStr []byte)(*Template,error){
	template := Template{}
	err := json.Unmarshal(templateStr,&template)
	if err != nil {
		return nil,err
	}
	template.init()
	return &template, nil
}