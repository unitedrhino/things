package info

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToSceneTypes(in *ud.SceneInfo) *types.SceneInfo {
	return &types.SceneInfo{
		ID:              in.Id,
		Name:            in.Name,
		Tag:             in.Tag,
		Desc:            in.Desc,
		If:              in.If,
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
		Tag:             in.Tag,
		Id:              in.ID,
		Name:            in.Name,
		Desc:            in.Desc,
		If:              in.If,
		When:            in.When,
		Then:            in.Then,
		Status:          in.Status,
		CreatedTime:     in.CreatedTime,
		HeadImg:         in.HeadImg,
		IsUpdateHeadImg: in.IsUpdateHeadImg,
	}
}
