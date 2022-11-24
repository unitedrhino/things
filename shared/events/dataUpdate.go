package events

type DataUpdateInfo struct {
	ProductID  string
	DeviceName string
	Data       any
}

type GatewayUpdateInfo struct {
	GatewayProductID  string
	GatewayDeviceName string
	Status            int32         //拓扑关系变化状态。* 2：解绑* 1：绑定
	Devices           []*DeviceCore //子设备列表
}

type DeviceCore struct {
	ProductID  string `json:"productID"`  //产品id
	DeviceName string `json:"deviceName"` //设备名称
}
