package devicegrouplogic

import (
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToGroupInfoPb(ro *relationDB.DmGroupInfo, productMap map[string]string) *dm.GroupInfo {
	if ro == nil {
		return nil
	}
	return &dm.GroupInfo{
		GroupID:     ro.GroupID,
		ParentID:    ro.ParentID,
		ProjectID:   ro.ProjectID,
		ProductName: productMap[ro.ProductID],
		GroupName:   ro.GroupName,
		ProductID:   ro.ProductID,
		Desc:        ro.Desc,
		CreatedTime: ro.CreatedTime.Unix(),
		Tags:        ro.Tags,
	}
}
