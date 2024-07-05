package userShared

const GetTypeAll = 1 //同时获取

type UserShareKey struct {
	ProductID    string `json:"productID"`  //产品id
	DeviceName   string `json:"deviceName"` //设备名称
	SharedUserID int64  `json:"sharedUserID"`
}
