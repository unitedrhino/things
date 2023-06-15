package utils

import (
	"github.com/i-Things/things/shared/def"
	"github.com/suifengtec/gocoord"
)

func PositionToBaidu(point def.Point) def.Point {
	gcp := gocoord.Position{Lon: point.Longitude, Lat: point.Latitude}
	switch point.CoordinateSystem {
	case def.CoordinateSystemMars:
		return ToDefPoint(def.CoordinateSystemBaidu, gocoord.GCJ02ToBD09(gcp))
	case def.CoordinateSystemEarth:
		return ToDefPoint(def.CoordinateSystemBaidu, gocoord.WGS84ToBD09(gcp))
	case def.CoordinateSystemBaidu:
		return point
	default:
		panic("暂未支持的坐标系")
	}
}
func PositionToEarth(point def.Point) def.Point {
	if point.CoordinateSystem == "" {
		point.CoordinateSystem = def.SystemCoordinate
	}
	gcp := gocoord.Position{Lon: point.Longitude, Lat: point.Latitude}
	switch point.CoordinateSystem {
	case def.CoordinateSystemMars:
		return ToDefPoint(def.CoordinateSystemEarth, gocoord.GCJ02ToWGS84(gcp))
	case def.CoordinateSystemEarth:
		return point
	case def.CoordinateSystemBaidu:
		return ToDefPoint(def.CoordinateSystemEarth, gocoord.BD09ToWGS84(gcp))
	default:
		panic("暂未支持的坐标系")
	}
}
func ToDefPoint(coordinateSystem def.CoordinateSystem, position gocoord.Position) def.Point {
	return def.Point{Longitude: position.Lon, Latitude: position.Lat, CoordinateSystem: coordinateSystem}
}
