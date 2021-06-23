package dict

type Profile struct {
	ProductID  string `json:"ProductId"`  //产品ID
	CategoryID string `json:"CategoryId"` //产品种类ID
}

type Maping struct {
	Zero string
	One  string
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
}

/*事件参数*/
type Param struct {
	ID     string `json:"id"`     //参数标识符
	Name   string `json:"name"`   //参数名称
	Define Define `json:"define"` //参数定义
}

/*事件*/
type Event struct {
	ID       string  `json:"id"`       //标识符
	Name     string  `json:"name"`     //功能名称
	Desc     string  `json:"desc"`     //描述
	Type     string  `json:"type"`     //事件类型: 信息:info  告警alert  故障:fault
	Params   []Param `json:"params"`   //事件参数
	Required bool    `json:"required"` //是否必须
}

/*行为*/
type Action struct {
	ID       string  `json:"id"`       //标识符
	Name     string  `json:"name"`     //功能名称
	Desc     string  `json:"desc"`     //描述
	Input    []Param `json:"input"`    //调用参数
	Output   []Param `json:"output"`   //返回参数
	Required bool    `json:"required"` //是否必须
}

/*属性*/
type Propert struct {
	ID       string `json:"id"`       //标识符
	Name     string `json:"name"`     //功能名称
	Desc     string `json:"gesc"`     //描述
	Mode     string `json:"mode"`     //读写乐行:rw(可读可写) r(只读)
	Define   Define `json:"define"`   //数据定义
	Required bool   `json:"required"` //是否必须
}

/*数据模板定义*/
type Template struct {
	Version    string    `json:"version"`    //版本号
	Properties []Propert `json:"properties"` //属性
	Events     []Event   `json:"events"`     //事件
	Actions    []Action  `json:"actions"`    //行为
	Profile    Profile   `json:"profile"`    //配置信息
}
