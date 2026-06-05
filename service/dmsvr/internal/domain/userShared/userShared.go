package userShared

import "time"

const GetTypeAll = 1 //同时获取
const MultiDeviceShareTokenTTL = 24 * time.Hour
const MultiDeviceShareTokenTTLSeconds = int64(MultiDeviceShareTokenTTL / time.Second)

type UserShareKey struct {
	ProductID    string `json:"productID"`  //产品id
	DeviceName   string `json:"deviceName"` //设备名称
	SharedUserID int64  `json:"sharedUserID"`
}
