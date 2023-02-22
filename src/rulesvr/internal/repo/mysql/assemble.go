package mysql

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
)

func ToScenePo(info *scene.Info) *RuleSceneInfo {
	ret := RuleSceneInfo{
		Id:          info.ID,
		Name:        info.Name,
		TriggerType: string(info.TriggerType),
		Trigger:     utils.AnyToNullString(info.Trigger),
		When:        utils.AnyToNullString(info.When),
		Then:        utils.AnyToNullString(info.Then),
		Desc:        info.Desc,
		State:       info.State,
	}
	switch info.TriggerType {
	case scene.TriggerTypeDevice:
		ret.Trigger = utils.AnyToNullString(info.Trigger.Device)
	}
	return &ret
}

func ToSceneDo(info *RuleSceneInfo) *scene.Info {
	ret := &scene.Info{
		ID:          info.Id,
		Name:        info.Name,
		Desc:        info.Desc,
		When:        make(scene.Terms, 0),
		Then:        make(scene.Actions, 0),
		State:       info.State,
		TriggerType: scene.TriggerType(info.TriggerType),
		CreatedTime: info.CreatedTime,
	}
	switch ret.TriggerType {
	case scene.TriggerTypeDevice:
		utils.SqlNullStringToAny(info.Trigger, &ret.Trigger.Device)
	}
	utils.SqlNullStringToAny(info.When, &ret.When)
	utils.SqlNullStringToAny(info.Then, &ret.Then)
	return ret
}
