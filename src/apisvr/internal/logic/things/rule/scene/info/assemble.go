package info

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
)

func ToSceneTypes(in *rule.SceneInfo) *types.SceneInfo {
	return &types.SceneInfo{
		ID:      in.Id,
		Name:    in.Name,
		Desc:    in.Desc,
		Trigger: in.Trigger,
		When:    in.When,
		Then:    in.Then,
	}
}
func ToScenePb(in *types.SceneInfo) *rule.SceneInfo {
	return &rule.SceneInfo{
		Id:      in.ID,
		Name:    in.Name,
		Desc:    in.Desc,
		Trigger: in.Trigger,
		When:    in.When,
		Then:    in.Then,
	}
}
