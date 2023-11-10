package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func transPbToApi(root *sys.AreaInfo) *types.AreaInfo {
	api := &types.AreaInfo{
		CreatedTime:  root.CreatedTime,
		ProjectID:    root.ProjectID,
		ParentAreaID: root.ParentAreaID,
		AreaID:       root.AreaID,
		AreaName:     root.AreaName,
		Position:     logic.ToSysPointApi(root.Position),
		Desc:         utils.ToNullString(root.Desc),
		Children:     nil,
	}
	if len(root.Children) > 0 {
		for _, child := range root.Children {
			api.Children = append(api.Children, transPbToApi(child))
		}
	}
	return api
}
