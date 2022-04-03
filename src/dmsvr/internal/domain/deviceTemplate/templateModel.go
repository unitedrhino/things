package deviceTemplate

//数据类型
type DataType string

const (
	BOOL      DataType = "bool"
	INT       DataType = "int"
	STRING    DataType = "string"
	STRUCT    DataType = "struct"
	FLOAT     DataType = "float"
	TIMESTAMP DataType = "timestamp"
	ARRAY     DataType = "array"
	ENUM      DataType = "enum"
)

type TEMP_TYPE int64

const (
	PROPERTY TEMP_TYPE = iota
	ACTION_INPUT
	ACTION_OUTPUT
	EVENT
)

type (
	/*数据模板定义*/
	Template struct {
		Version    string               `json:"version"`    //版本号
		Properties []Property           `json:"properties"` //属性
		Events     []Event              `json:"events"`     //事件
		Actions    []Action             `json:"actions"`    //行为
		Profile    Profile              `json:"profile"`    //配置信息
		Property   map[string]*Property `json:"-"`          //内部使用,使用map加速匹配,key为id
		Event      map[string]*Event    `json:"-"`          //内部使用,使用map加速匹配,key为id
		Action     map[string]*Action   `json:"-"`          //内部使用,使用map加速匹配,key为id
	}
	Profile struct {
		ProductID  string `json:"ProductId"`  //产品ID
		CategoryID string `json:"CategoryId"` //产品种类ID
	}
	/*结构体说明*/
	Spec struct {
		ID       string `json:"id"`       //参数标识符
		Name     string `json:"name"`     //参数名称
		DataType Define `json:"dataType"` //参数定义
	}

	/*参数*/
	Param struct {
		ID     string `json:"id"`               //参数标识符
		Name   string `json:"name"`             //参数名称
		Define Define `json:"define,omitempty"` //参数定义
	}

	/*事件*/
	Event struct {
		ID       string            `json:"id"`       //标识符
		Name     string            `json:"name"`     //功能名称
		Desc     string            `json:"desc"`     //描述
		Type     string            `json:"type"`     //事件类型: 信息:info  告警alert  故障:fault
		Params   []Param           `json:"params"`   //事件参数
		Required bool              `json:"required"` //是否必须
		Param    map[string]*Param `json:"-"`        //内部使用,使用map加速匹配,key为id
	}
	/*行为*/
	Action struct {
		ID       string            `json:"id"`       //标识符
		Name     string            `json:"name"`     //功能名称
		Desc     string            `json:"desc"`     //描述
		Input    []Param           `json:"input"`    //调用参数
		Output   []Param           `json:"output"`   //返回参数
		Required bool              `json:"required"` //是否必须
		In       map[string]*Param `json:"-"`        //内部使用,使用map加速匹配,key为id
		Out      map[string]*Param `json:"-"`        //内部使用,使用map加速匹配,key为id
	}
	/*属性*/
	Property struct {
		ID       string `json:"id"`       //标识符
		Name     string `json:"name"`     //功能名称
		Desc     string `json:"gesc"`     //描述
		Mode     string `json:"mode"`     //读写乐行:rw(可读可写) r(只读)
		Define   Define `json:"define"`   //数据定义
		Required bool   `json:"required"` //是否必须
	}
	/*数据类型定义*/
	Define struct {
		Type      DataType          `json:"type"`                //参数类型:bool int string struct float timestamp array enum
		Maping    map[string]string `json:"mapping,omitempty"`   //枚举及bool类型:bool enum
		Min       string            `json:"min,omitempty"`       //数值最小值:int  float
		Max       string            `json:"max,omitempty"`       //数值最大值:int string float
		Start     string            `json:"start,omitempty"`     //初始值:int float
		Step      string            `json:"step,omitempty"`      //步长:int float
		Unit      string            `json:"unit,omitempty"`      //单位:int float
		Specs     []Spec            `json:"specs,omitempty"`     //结构体:struct
		ArrayInfo *Define           `json:"arrayInfo,omitempty"` //数组:array
		Spec      map[string]*Spec  //内部使用,使用map加速匹配,key为id
	}
)
