package def

import (
	"github.com/i-Things/things/shared/store"
)

const DeviceGroupLevel = 3

type CoordinateSystem string

var SystemCoordinate = CoordinateSystemBaidu //默认坐标系

const (
	CoordinateSystemEarth CoordinateSystem = "WGS84" //GPS坐标系：地球系
	CoordinateSystemMars  CoordinateSystem = "GCJ02" //GPS坐标系：火星系
	CoordinateSystemBaidu CoordinateSystem = "BD09"  //GPS坐标系：百度系
)

// 坐标，
type Point struct {
	CoordinateSystem CoordinateSystem `json:"coordinateSystem,omitempty"` //坐标系：WGS84(地球系)，GCJ02(火星系)，BD09(百度系)<br/>参考解释：https://www.cnblogs.com/bigroc/p/16423120.html
	Longitude        float64          `json:"longitude,range=[0:180]"`    //经度
	Latitude         float64          `json:"latitude,range=[0:90]"`      //纬度
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
