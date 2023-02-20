// Package scene 执行动作
package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type ActionAlarmMode string

const (
	ActionAlarmModeTrigger ActionAlarmMode = "trigger"
	ActionAlarmModeRelieve ActionAlarmMode = "relieve"
)

type ActionAlarm struct {
	Mode ActionAlarmMode `json:"mode"` //告警模式  trigger: 触发告警  relieve: 解除告警
}

func (a *ActionAlarm) Validate() error {
	if a == nil {
		return nil
	}
	if !utils.SliceIn(a.Mode, ActionAlarmModeRelieve, ActionAlarmModeTrigger) {
		return errors.Parameter.AddMsg("操作告警不支持的类型:" + string(a.Mode))
	}
	return nil
}
