package info

import (
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToDeviceInfoCount(in *dm.DeviceInfoCount) *types.DeviceInfoCount {
	if in == nil {
		return nil
	}
	return &types.DeviceInfoCount{
		Total:    in.Total,
		Online:   in.Online,
		Offline:  in.Offline,
		Inactive: in.Inactive,
		Unknown:  in.Unknown,
	}
}

func ToGroupInfoTypes(in *dm.GroupInfo) *types.GroupInfo {
	return &types.GroupInfo{
		AreaID:          in.AreaID,
		ProductID:       in.ProductID,
		ProductName:     in.ProductName,
		ID:              in.Id,
		ParentID:        in.ParentID,
		ProjectID:       in.ProjectID,
		Name:            in.Name,
		CreatedTime:     in.CreatedTime,
		Desc:            in.Desc,
		Tags:            logic.ToTagsType(in.Tags),
		DeviceInfoCount: ToDeviceInfoCount(in.DeviceInfoCount),
	}
}
func ToGroupInfoPbTypes(in *types.GroupInfo) *dm.GroupInfo {
	return &dm.GroupInfo{
		AreaID:      in.AreaID,
		ProductID:   in.ProductID,
		ProductName: in.ProductName,
		Id:          in.ID,
		ParentID:    in.ParentID,
		ProjectID:   in.ProjectID,
		Name:        in.Name,
		CreatedTime: in.CreatedTime,
		Desc:        in.Desc,
		Tags:        logic.ToTagsMap(in.Tags),
	}
}
