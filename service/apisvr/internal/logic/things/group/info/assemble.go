package info

import (
	"gitee.com/i-Things/things/service/apisvr/internal/logic"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToGroupInfosTypes(in []*dm.GroupInfo) []*types.GroupInfo {
	if len(in) == 0 {
		return nil
	}
	glist := make([]*types.GroupInfo, 0, len(in))
	for _, v := range in {
		glist = append(glist, ToGroupInfoTypes(v))
	}
	return glist
}

func ToGroupInfoTypes(in *dm.GroupInfo) *types.GroupInfo {
	return &types.GroupInfo{
		AreaID:      in.AreaID,
		ProductID:   in.ProductID,
		ProductName: in.ProductName,
		ID:          in.Id,
		ParentID:    in.ParentID,
		ProjectID:   in.ProjectID,
		Name:        in.Name,
		CreatedTime: in.CreatedTime,
		Desc:        in.Desc,
		DeviceCount: in.DeviceCount,
		IsLeaf:      in.IsLeaf,
		Tags:        logic.ToTagsType(in.Tags),
		Children:    ToGroupInfosTypes(in.Children),
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
