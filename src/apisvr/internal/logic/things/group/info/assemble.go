package info

import (
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToGroupInfoTypes(in *dm.GroupInfo) *types.GroupInfo {
	return &types.GroupInfo{
		ProductID:   in.ProductID,
		ProductName: in.ProductName,
		GroupID:     in.GroupID,
		ParentID:    in.ParentID,
		ProjectID:   in.ProjectID,
		GroupName:   in.GroupName,
		CreatedTime: in.CreatedTime,
		Desc:        in.Desc,
		Tags:        logic.ToTagsType(in.Tags),
	}
}
