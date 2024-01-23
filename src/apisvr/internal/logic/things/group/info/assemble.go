package info

import (
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

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
		Tags:        logic.ToTagsType(in.Tags),
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
