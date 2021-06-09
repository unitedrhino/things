package model

type Profile struct {
	ProductID  string //产品ID
	CategoryID string //产品种类ID
}

//type Maping struct {
//	Zero string
//	One  string
//}

/*结构体说明*/
type Spec struct {
	ID       string //参数标识符
	Name     string //参数名称
	DataType Define //参数定义
}

/*数据类型定义*/
type Define struct {
	Type      string            //参数类型:bool int string struct float timestamp array enum
	Maping    map[string]string //枚举及bool类型:bool enum
	Min       string            //数值最小值:int string float
	Max       string            //数值最大值:int string float
	Start     string            //初始值:int float
	Step      string            //步长:int float
	Unit      string            //单位:int float
	Specs     *Spec             //结构体:struct
	ArrayInfo *Define           //数组:array
}

/*事件参数*/
type Param struct {
	ID     string //参数标识符
	Name   string //参数名称
	Define Define //参数定义
}

/*事件*/
type Event struct {
	ID       string  //标识符
	Name     string  //功能名称
	Desc     string  //描述
	Type     string  //事件类型: 信息:info  告警alert  故障:fault
	Params   []Param //事件参数
	Required bool    //是否必须
}

/*行为*/
type Action struct {
	ID       string  //标识符
	Name     string  //功能名称
	Desc     string  //描述
	Input    []Param //调用参数
	Output   []Param //返回参数
	Required bool    //是否必须
}

/*属性*/
type Propert struct {
	ID       string   //标识符
	Name     string   //功能名称
	Desc     string   //描述
	Mode     string   //读写乐行:rw(可读可写) r(只读)
	Define   []Define //数据定义
	Required bool     //是否必须
}

type Template struct {
	Version    string  //版本号
	Properties Propert //属性
	Events     Event   //事件
	Actions    Action  //行为
	Profile    Profile //配置信息
}
