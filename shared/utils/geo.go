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
	default: //默认百度坐标系
		return point
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
	default: //默认百度坐标系
		return ToDefPoint(def.CoordinateSystemEarth, gocoord.BD09ToWGS84(gcp))
	}
}
func ToDefPoint(coordinateSystem def.CoordinateSystem, position gocoord.Position) def.Point {
	return def.Point{Longitude: position.Lon, Latitude: position.Lat, CoordinateSystem: coordinateSystem}
}
