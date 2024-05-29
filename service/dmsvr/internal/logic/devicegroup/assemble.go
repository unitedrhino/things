package devicegrouplogic

import (
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToGroupInfoPb(ro *relationDB.DmGroupInfo) *dm.GroupInfo {
	if ro == nil {
		return nil
	}
	productName := ""
	if ro.ProductInfo != nil {
		productName = ro.ProductInfo.ProductName
	}
	return &dm.GroupInfo{
		AreaID:      int64(ro.AreaID),
		Id:          ro.ID,
		ParentID:    ro.ParentID,
		ProjectID:   int64(ro.ProjectID),
		ProductName: productName,
		Name:        ro.Name,
		ProductID:   ro.ProductID,
		DeviceCount: ro.DeviceCount,
		Desc:        ro.Desc,
		CreatedTime: ro.CreatedTime.Unix(),
		Tags:        ro.Tags,
		IsLeaf:      ro.IsLeaf,
	}
}
