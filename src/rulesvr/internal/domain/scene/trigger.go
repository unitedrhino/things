// Package scene 触发器
package scene

import (
	"github.com/i-Things/things/shared/errors"
)

type Trigger struct {
	Device TriggerDevices `json:"device"`
	Timer  *Timer         `json:"timer"`
}

func (t *Trigger) Validate(triggerType TriggerType) error {
	if t == nil {
		return nil
	}
	switch triggerType {
	case TriggerTypeDevice:
		if t.Device == nil {
			return errors.Parameter.AddMsg("设备类型触发器需要填写触发详情")
		}
		return t.Device.Validate()
	}
	return nil
}
