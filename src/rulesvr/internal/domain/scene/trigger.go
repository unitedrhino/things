// Package scene 触发器
package scene

type Trigger struct {
	Type   string         `json:"type"` //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Device *TriggerDevice `json:"device"`
}
type TriggerDevice struct {
	ProductID      string          `json:"productID"`      //产品id
	Selector       string          `json:"selector"`       //设备选择方式  all: 全部 fixed:指定的设备
	SelectorValues []string        `json:"selectorValues"` //选择的列表  选择的列表, fixed类型是设备名列表
	Operation      DeviceOperation `json:"operation"`
}
type DeviceOperation struct {
	Operator string `json:"operator"` //触发类型  online:上线 offline:下线 reportProperty:属性上报 reportEvent: 事件上报
}
