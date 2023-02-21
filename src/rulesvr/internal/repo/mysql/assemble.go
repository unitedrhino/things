package mysql

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
)

func ToScenePo(info *scene.Info) *RuleSceneInfo {
	return &RuleSceneInfo{
		Id:      info.ID,
		Name:    info.Name,
		Trigger: utils.AnyToNullString(info.Trigger),
		When:    utils.AnyToNullString(info.When),
		Then:    utils.AnyToNullString(info.Then),
		Desc:    info.Desc,
	}
}

func ToSceneDo(info *RuleSceneInfo) *scene.Info {
	ret := &scene.Info{
		ID:   info.Id,
		Name: info.Name,
		Desc: info.Desc,
		When: make(scene.Terms, 0),
		Then: make(scene.Actions, 0),
	}
	utils.SqlNullStringToAny(info.Trigger, &ret.Trigger)
	utils.SqlNullStringToAny(info.When, &ret.When)
	utils.SqlNullStringToAny(info.Then, &ret.Then)
	return ret
}
