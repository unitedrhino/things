package info

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/udsvr/pb/ud"
)

func ToSceneTypes(in *ud.SceneInfo) *types.SceneInfo {
	return &types.SceneInfo{
		ID:              in.Id,
		Name:            in.Name,
		AreaIDs:         in.AreaIDs,
		Tag:             in.Tag,
		Desc:            in.Desc,
		Trigger:         in.Trigger,
		When:            in.When,
		Then:            in.Then,
		Status:          in.Status,
		CreatedTime:     in.CreatedTime,
		HeadImg:         in.HeadImg,
		IsUpdateHeadImg: in.IsUpdateHeadImg,
	}
}
func ToScenePb(in *types.SceneInfo) *ud.SceneInfo {
	return &ud.SceneInfo{
		AreaIDs:         in.AreaIDs,
		Tag:             in.Tag,
		Id:              in.ID,
		Name:            in.Name,
		Desc:            in.Desc,
		Trigger:         in.Trigger,
		When:            in.When,
		Then:            in.Then,
		Status:          in.Status,
		CreatedTime:     in.CreatedTime,
		HeadImg:         in.HeadImg,
		IsUpdateHeadImg: in.IsUpdateHeadImg,
	}
}
