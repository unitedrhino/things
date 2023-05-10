package utils

import (
	"github.com/i-Things/things/shared/def"
	"github.com/suifengtec/gocoord"
)

func PositionToBaidu(cs def.CoordinateSystem, lon, lat float64) gocoord.Position {
	gcp := gocoord.Position{lon, lat}
	switch cs {
	case def.CoordinateSystemMars:
		return gocoord.GCJ02ToBD09(gcp)
	case def.CoordinateSystemEarth:
		return gocoord.WGS84ToBD09(gcp)
	case def.CoordinateSystemBaidu:
		return gcp
	default:
		panic("暂未支持的坐标系")
	}
}
