package relationDB

import (
	"github.com/i-Things/things/service/rulesvr/internal/domain/scene"
)

func SceneInfoDoToPo(info *scene.Info) *RuleSceneInfo {
	return &RuleSceneInfo{
		ID:          info.ID,
		Name:        info.Name,
		TriggerType: string(info.TriggerType),
		Trigger:     info.Trigger,
		When:        info.When,
		Then:        info.Then,
		Desc:        info.Desc,
		Status:      info.Status,
	}
}
func SceneInfoPoToDo(info *RuleSceneInfo) *scene.Info {
	return &scene.Info{
		ID:          info.ID,
		Name:        info.Name,
		TriggerType: scene.TriggerType(info.TriggerType),
		Trigger:     info.Trigger,
		When:        info.When,
		Then:        info.Then,
		Desc:        info.Desc,
		Status:      info.Status,
	}
}

func SceneInfoPoToDos(info []*RuleSceneInfo) (ret scene.Infos) {
	if info == nil {
		return
	}
	for _, v := range info {
		ret = append(ret, SceneInfoPoToDo(v))
	}
	return
}
