package devices

type Core struct {
	ProductID  string `json:"productID"`  //产品id
	DeviceName string `json:"deviceName"` //设备名称
}

// 设备标签
type Tag struct {
	Key   string `json:"key"`   //设备标签key
	Value string `json:"value"` //设备标签value
}
