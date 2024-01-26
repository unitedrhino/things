package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToAreaInfoTypes(root *sys.AreaInfo, deviceCount *types.DeviceInfoCount) *types.AreaInfo {
	api := &types.AreaInfo{
		CreatedTime:     root.CreatedTime,
		ProjectID:       root.ProjectID,
		ParentAreaID:    root.ParentAreaID,
		AreaID:          root.AreaID,
		AreaName:        root.AreaName,
		AreaNamePath:    root.AreaNamePath,
		LowerLevelCount: root.LowerLevelCount,
		AreaIDPath:      root.AreaIDPath,
		Position:        logic.ToSysPointApi(root.Position),
		Desc:            utils.ToNullString(root.Desc),
		DeviceInfoCount: deviceCount,
		Children:        nil,
	}
	if len(root.Children) > 0 {
		for _, child := range root.Children {
			api.Children = append(api.Children, ToAreaInfoTypes(child, nil))
		}
	}
	return api
}
func ToAreaInfosTypes(in []*sys.AreaInfo, deviceCounts []*types.DeviceInfoCount) (ret []*types.AreaInfo) {
	for i, v := range in {
		ret = append(ret, ToAreaInfoTypes(v, deviceCounts[i]))
	}
	return
}
