// Package scene 执行动作
package scene

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type ActionAlarmMode = string

const (
	ActionAlarmModeTrigger ActionAlarmMode = "trigger" //触发告警
	ActionAlarmModeRelieve ActionAlarmMode = "relieve" //接触告警
)

type AlarmDealStatus = int64

const (
	AlarmDealStatusWaring    = 1 //告警中
	AlarmDealStatusIgnore    = 2 //忽略
	AlarmDealStatusInHand    = 3 //正在处理
	AlarmDealStatusProcessed = 4 //已处理
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
	status := int64(def.True)
	er := errors.Fmt(err)
	if er.GetCode() != errors.OK.GetCode() {
		status = def.False
		repo.Info.Log.Status = def.False
	}
	repo.Info.Log.ActionMutex.Lock()
	defer repo.Info.Log.ActionMutex.Unlock()
	repo.Info.Log.Actions = append(repo.Info.Log.Actions, &LogAction{
		Type: ActionExecutorAlarm,
		Alarm: &LogActionAlarm{
			Mode: a.Mode,
		},
		Status: status,
		Code:   er.GetCode(),
		Msg:    er.GetMsg(),
	})
	return err
}
