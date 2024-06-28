package logic

import (
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToPageInfo(info *dm.PageInfo) *stores.PageInfo {
	return utils.Copy[stores.PageInfo](info)
}

func ToDmPoint(point *stores.Point) *dm.Point {
	if point == nil {
		return nil
	}
	return &dm.Point{Longitude: point.Longitude, Latitude: point.Latitude}
}
func ToStorePoint(point *dm.Point) stores.Point {
	if point == nil {
		return stores.Point{Longitude: 0, Latitude: 0}
	}
	return stores.Point{Longitude: point.Longitude, Latitude: point.Latitude}
}
func ToDeviceCores(in []*dm.DeviceCore) []*devices.Core {
	list := make([]*devices.Core, 0, len(in))
	for _, v := range in {
		if v == nil {
			continue
		}
		list = append(list, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return list
}
