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

type AlarmTrigger struct {
	SceneID     int64  //场景ID
	SceneName   string //场景名称
	TriggerType int64  //触发类型(设备触发1,其他2)
	ProductID   string //触发产品id
	DeviceName  string //触发设备名称
	Serial      string //告警流水
	Desc        string //告警说明
}

type AlarmRelieve struct {
	SceneID   int64  //场景ID
	SceneName string //场景名称
}

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
