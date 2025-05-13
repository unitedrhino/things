package devices

const (
	DeviceRegisterUnable   = 1
	DeviceAutoCreateEnable = 3
	EncryptionTypeCert     = 1
)

type Core struct {
	ProductID  string `json:"productID"`  //产品id
	DeviceName string `json:"deviceName"` //设备名称
}

type WithGateway struct {
	Dev     Core  `json:"dev"`
	Gateway *Core `json:"gateway,omitempty"` //如果是子设备类型,会带上网关
}

type Info struct {
	ProductID    string `json:"productID"`  //产品id
	DeviceName   string `json:"deviceName"` //设备名称
	TenantCode   string
	ProjectID    int64
	AreaID       int64
	AreaIDPath   string
	GroupIDs     []int64
	GroupIDPaths []string
}

// 归属
type Affiliation struct {
	TenantCode   string
	ProjectID    int64
	AreaID       int64
	AreaIDPath   string
	GroupIDs     []int64
	GroupIDPaths []string
}

// 设备标签
type Tag struct {
	Key   string `json:"key"`   //设备标签key
	Value string `json:"value"` //设备标签value
}
