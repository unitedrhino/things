package scene

type ActionDevice struct {
	ProductID      string                `json:"productID"`      //产品id
	Selector       TriggerDeviceSelector `json:"selector"`       //设备选择方式   fixed:指定的设备
	SelectorValues []string              `json:"selectorValues"` //选择的列表  选择的列表, fixed类型是设备名列表
	Type           string                `json:"type"`           // 云端向设备发起属性控制: propertyControl  应用调用设备行为:action  todo:通知设备上报
	DataID         []string              `json:"dataID"`         // 属性的id及事件的id aa.bb.cc
	Value          string                `json:"value"`          //传的值
}
