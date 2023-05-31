package devices

import (
	"github.com/i-Things/things/shared/store"
)

type Core struct {
	ProductID  string `json:"productID"`  //产品id
	DeviceName string `json:"deviceName"` //设备名称
}

// 设备标签
type Tag struct {
	Key   string `json:"key"`   //设备标签key
	Value string `json:"value"` //设备标签value
}

// 设备位置坐标，
type Point struct {
	Longitude float64 `json:"longitude,range=[0:180]"` //经度
	Latitude  float64 `json:"latitude,range=[0:90]"`   //纬度
}

func (p Point) ToPo() store.Point {
	return store.Point{
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}
}
func ToPoint(p store.Point) Point {
	return Point{
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}
}
