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
		Trigger: ToSceneTriggerTypes(in.Trigger),
		When:    ToSceneTermsTypes(in.When),
		Then:    ToSceneActionsTypes(in.Then),
	}
}

func ToSceneActionsTypes(in []*rule.SceneAction) []*types.SceneAction {
	if in == nil {
		return nil
	}
	actions := make([]*types.SceneAction, 0, len(in))
	for _, v := range in {
		actions = append(actions, ToSceneActionTypes(v))
	}
	return actions
}

func ToSceneActionTypes(in *rule.SceneAction) *types.SceneAction {
	return &types.SceneAction{Executor: in.Executor}
}

func ToSceneActionDelayTypes(in *rule.SceneActionDelay) *types.SceneActionDelay {
	return &types.SceneActionDelay{
		Time: in.Time,
		Unit: in.Unit,
	}
}
func ToSceneActionAlarmTypes(in *rule.SceneActionAlarm) *types.SceneActionAlarm {
	return &types.SceneActionAlarm{Mode: in.Mode}
}

func ToSceneTermsTypes(in []*rule.SceneTerm) []*types.SceneTerm {
	if in == nil {
		return nil
	}
	terms := make([]*types.SceneTerm, 0, len(in))
	for _, v := range in {
		terms = append(terms, ToSceneTermTypes(v))
	}
	return terms
}

func ToSceneTermTypes(in *rule.SceneTerm) *types.SceneTerm {
	return &types.SceneTerm{
		Column:   in.Column,
		Value:    in.Value,
		Type:     in.Type,
		TermType: in.TermType,
		Terms:    ToSceneTermsTypes(in.Terms),
	}
}

func ToSceneTriggerTypes(in *rule.SceneTrigger) *types.SceneTrigger {
	return &types.SceneTrigger{Type: in.Type, Device: ToSceneTriggerDeviceTypes(in.Device)}
}
func ToSceneTriggerDeviceTypes(in *rule.SceneTriggerDevice) *types.SceneTriggerDevice {
	return &types.SceneTriggerDevice{
		ProductID:      in.ProductID,
		Selector:       in.Selector,
		SelectorValues: in.SelectorValues,
		Operation:      types.SceneDeviceOperation{Operator: in.Operation.Operator},
	}
}
