package info

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
)

func ToScenePb(in *types.SceneInfo) *rule.SceneInfo {
	return &rule.SceneInfo{
		Id:      in.ID,
		Name:    in.Name,
		Desc:    in.Desc,
		Trigger: ToSceneTriggerPb(in.Trigger),
		When:    ToSceneTermsPb(in.When),
		Then:    ToSceneActionsPb(in.Then),
	}
}

func ToSceneActionsPb(in []*types.SceneAction) []*rule.SceneAction {
	if in == nil {
		return nil
	}
	if in == nil {
		return nil
	}
	actions := make([]*rule.SceneAction, 0, len(in))
	for _, v := range in {
		actions = append(actions, ToSceneActionPb(v))
	}
	return actions
}

func ToSceneActionPb(in *types.SceneAction) *rule.SceneAction {
	if in == nil {
		return nil
	}
	return &rule.SceneAction{Executor: in.Executor}
}

func ToSceneActionDelayPb(in *types.SceneActionDelay) *rule.SceneActionDelay {
	if in == nil {
		return nil
	}
	return &rule.SceneActionDelay{
		Time: in.Time,
		Unit: in.Unit,
	}
}
func ToSceneActionAlarmPb(in *types.SceneActionAlarm) *rule.SceneActionAlarm {
	if in == nil {
		return nil
	}
	return &rule.SceneActionAlarm{Mode: in.Mode}
}

func ToSceneTermsPb(in []*types.SceneTerm) []*rule.SceneTerm {
	if in == nil {
		return nil
	}
	terms := make([]*rule.SceneTerm, 0, len(in))
	for _, v := range in {
		terms = append(terms, ToSceneTermPb(v))
	}
	return terms
}

func ToSceneTermPb(in *types.SceneTerm) *rule.SceneTerm {
	if in == nil {
		return nil
	}
	return &rule.SceneTerm{
		Column:   in.Column,
		Value:    in.Value,
		Type:     in.Type,
		TermType: in.TermType,
		Terms:    ToSceneTermsPb(in.Terms),
	}
}

func ToSceneTriggerPb(in *types.SceneTrigger) *rule.SceneTrigger {
	if in == nil {
		return nil
	}
	return &rule.SceneTrigger{Type: in.Type, Device: ToSceneTriggerDevicePb(in.Device)}
}
func ToSceneTriggerDevicePb(in *types.SceneTriggerDevice) *rule.SceneTriggerDevice {
	if in == nil {
		return nil
	}
	return &rule.SceneTriggerDevice{
		ProductID:      in.ProductID,
		Selector:       in.Selector,
		SelectorValues: in.SelectorValues,
		Operation:      &rule.SceneTriggerDeviceOperation{Operator: in.Operation.Operator},
	}
}
