// Package scene 执行动作
package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type ActionAlarmMode string

const (
	ActionAlarmModeTrigger ActionAlarmMode = "trigger" //触发告警
	ActionAlarmModeRelieve ActionAlarmMode = "relieve" //接触告警
)

type ActionAlarm struct {
	Mode ActionAlarmMode `json:"mode"` //告警模式  trigger: 触发告警  relieve: 解除告警
}

type AlarmSerial struct {
	Scene *Info           //场景ID
	Mode  ActionAlarmMode //告警模式
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

func (a *ActionAlarm) Execute(ctx context.Context, repo ActionRepo) error {
	err := repo.AlarmExec(ctx, AlarmSerial{Mode: a.Mode, Scene: repo.Info})
	return err
}
