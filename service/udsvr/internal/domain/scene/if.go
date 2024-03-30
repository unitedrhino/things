// Package scene 触发器
package scene

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
)

type If struct {
	Triggers Triggers `json:"triggers,omitempty"`
}

type TriggerType = string

const (
	TriggerTypeDevice TriggerType = "device"
	TriggerTypeTimer  TriggerType = "timer"
)

type Triggers []*Trigger
type Trigger struct {
	Type   TriggerType    `json:"type"`
	Order  int64          `json:"order"`
	Device *TriggerDevice `json:"device,omitempty"` //设备触发
	Timer  *Timer         `json:"timer,omitempty"`  //定时触发
}

func (t Triggers) Validate(repo ValidateRepo) error {
	if len(t) == 0 {
		return nil
	}
	for _, v := range t {
		err := v.Validate(repo)
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *Trigger) Validate(repo ValidateRepo) error {
	if t == nil {
		return errors.Parameter.AddMsg("需要填写触发内容")
	}
	if !utils.SliceIn(t.Type, TriggerTypeTimer, TriggerTypeDevice) {
		return errors.Parameter.AddMsg("触发器不支持的类型:" + string(t.Type))
	}
	switch t.Type {
	case TriggerTypeDevice:
		return t.Device.Validate(repo)
	case TriggerTypeTimer:
		return t.Timer.Validate()
	}
	return nil
}

func (i *If) Validate(t SceneType, repo ValidateRepo) error {
	switch t {
	case SceneTypeManual:
		return nil
	case SceneTypeAuto:
		if len(i.Triggers) == 0 {
			return errors.Parameter.AddMsg("自动触发类型需要填写至少一项设备触发或者定时触发")
		}
		err := i.Triggers.Validate(repo)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.Parameter.AddMsg("触发类型只支持: auto:自动触发 manual:手动触发")
	}
}
