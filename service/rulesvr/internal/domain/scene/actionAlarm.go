// Package scene 执行动作
package scene

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type ActionAlarmMode string

const (
	ActionAlarmModeTrigger ActionAlarmMode = "trigger"
	ActionAlarmModeRelieve ActionAlarmMode = "relieve"
)

type TriggerSerial struct {
	SceneID     int64  //场景ID
	SceneName   string //场景名称
	TriggerType int64  //触发类型(设备触发1,其他2)
	Device      devices.Core
	Serial      string //流水
	Desc        string //说明
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
func (a *ActionAlarm) Execute(ctx context.Context, repo ActionRepo) error {
	switch a.Mode {
	case ActionAlarmModeRelieve:
		err := repo.Alarm.AlarmRelieve(ctx, AlarmRelieve{
			SceneID:   repo.Scene.ID,
			SceneName: repo.Scene.Name,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.AlarmRelieve err:%v", utils.FuncName(), err)
			return err
		}
	case ActionAlarmModeTrigger:
		serial := ""
		if repo.Serial != nil {
			serial = repo.Serial.GenSerial()
		}
		err := repo.Alarm.AlarmTrigger(ctx, TriggerSerial{
			SceneID:     repo.Scene.ID,
			SceneName:   repo.Scene.Name,
			TriggerType: 1,
			Device:      repo.Device,
			Serial:      serial,
			Desc:        "",
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.AlarmTrigger err:%v", utils.FuncName(), err)
			return err
		}
	}
	return nil
}
