// Package scene 触发器
package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type TriggerType string

const (
	TriggerTypeDevice TriggerType = "device"
	TriggerTypeTimer  TriggerType = "timer"
	TriggerTypeManual TriggerType = "manual"
)

type Trigger struct {
	Type   TriggerType    `json:"type"` //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Device *TriggerDevice `json:"device"`
}

func (t *Trigger) Validate() error {
	if t == nil {
		return nil
	}
	if !utils.SliceIn(t.Type, TriggerTypeDevice) {
		return errors.Parameter.AddMsg("触发器不支持的类型:" + string(t.Type))
	}
	if t.Device == nil {
		return errors.Parameter.AddMsg("设备类型触发器需要填写触发详情")
	}
	return t.Device.Validate()
}
