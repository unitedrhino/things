package devicegrouplogic

import (
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
		GroupID:     ro.GroupID,
		ParentID:    ro.ParentID,
		ProjectID:   ro.ProjectID,
		ProductName: productName,
		GroupName:   ro.GroupName,
		ProductID:   ro.ProductID,
		Desc:        ro.Desc,
		CreatedTime: ro.CreatedTime.Unix(),
		Tags:        ro.Tags,
	}
}
