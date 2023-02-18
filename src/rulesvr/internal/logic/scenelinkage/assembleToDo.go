package scenelinkagelogic

import (
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
)

func ToSceneDo(in *rule.SceneInfo) *scene.Info {
	return &scene.Info{
		ID:      in.Id,
		Name:    in.Name,
		Desc:    in.Desc,
		Trigger: ToSceneTriggerDo(in.Trigger),
		When:    ToSceneTermsDo(in.When),
		Then:    ToSceneActionsDo(in.Then),
	}
}

func ToSceneActionsDo(in []*rule.SceneAction) []*scene.Action {
	if in == nil {
		return nil
	}
	if in == nil {
		return nil
	}
	actions := make([]*scene.Action, 0, len(in))
	for _, v := range in {
		actions = append(actions, ToSceneActionDo(v))
	}
	return actions
}

func ToSceneActionDo(in *rule.SceneAction) *scene.Action {
	if in == nil {
		return nil
	}
	return &scene.Action{Executor: in.Executor}
}

func ToSceneActionDelayDo(in *rule.SceneActionDelay) *scene.ActionDelay {
	if in == nil {
		return nil
	}
	return &scene.ActionDelay{
		Time: in.Time,
		Unit: in.Unit,
	}
}
func ToSceneActionAlarmDo(in *rule.SceneActionAlarm) *scene.ActionAlarm {
	if in == nil {
		return nil
	}
	return &scene.ActionAlarm{Mode: in.Mode}
}

func ToSceneTermsDo(in []*rule.SceneTerm) []*scene.Term {
	if in == nil {
		return nil
	}
	terms := make([]*scene.Term, 0, len(in))
	for _, v := range in {
		terms = append(terms, ToSceneTermDo(v))
	}
	return terms
}

func ToSceneTermDo(in *rule.SceneTerm) *scene.Term {
	if in == nil {
		return nil
	}
	return &scene.Term{
		Column:   in.Column,
		Value:    in.Value,
		Type:     in.Type,
		TermType: in.TermType,
		Terms:    ToSceneTermsDo(in.Terms),
	}
}

func ToSceneTriggerDo(in *rule.SceneTrigger) *scene.Trigger {
	if in == nil {
		return nil
	}
	return &scene.Trigger{Type: in.Type, Device: ToSceneTriggerDeviceDo(in.Device)}
}
func ToSceneTriggerDeviceDo(in *rule.SceneTriggerDevice) *scene.TriggerDevice {
	if in == nil {
		return nil
	}
	return &scene.TriggerDevice{
		ProductID:      in.ProductID,
		Selector:       in.Selector,
		SelectorValues: in.SelectorValues,
		Operation:      scene.DeviceOperation{Operator: in.Operation.Operator},
	}
}
