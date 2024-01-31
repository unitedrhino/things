// Package scene 触发器
package scene

import (
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
)

type TriggerType string

const (
	//TriggerTypeDevice TriggerType = "device"
	//TriggerTypeTimer  TriggerType = "timer"
	TriggerTypeManual TriggerType = "manual"
	TriggerTypeAuto   TriggerType = "auto"
)

type Trigger struct {
	Type    TriggerType    `json:"type"`    //触发类型 auto: 自动触发 manual:手动触发
	Devices TriggerDevices `json:"devices"` //设备触发
	Timers  Timers         `json:"timers"`  //定时触发
}

func (t *Trigger) Validate() error {
	if t == nil {
		return errors.Parameter.AddMsg("需要填写触发内容")
	}
	if !utils.SliceIn(t.Type, TriggerTypeAuto, TriggerTypeManual) {
		return errors.Parameter.AddMsg("触发器不支持的类型:" + string(t.Type))
	}
	switch t.Type {
	case TriggerTypeManual:
		return nil
	case TriggerTypeAuto:
		if len(t.Devices) == 0 && len(t.Timers) == 0 {
			return errors.Parameter.AddMsg("自动触发类型需要填写至少一项设备触发或者定时触发")
		}
		err := t.Timers.Validate()
		if err != nil {
			return err
		}
		err = t.Devices.Validate()
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.Parameter.AddMsg("触发类型只支持: auto:自动触发 manual:手动触发")
	}
}
