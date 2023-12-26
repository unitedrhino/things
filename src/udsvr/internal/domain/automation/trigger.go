// Package scene 触发器
package automation

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
	case TriggerTypeManual:
		return nil
	case TriggerTypeTimer:
		if t.Timer == nil {
			return errors.Parameter.AddMsg("时间类型触发器需要填写触发详情")
		}
		return t.Timer.Validate()
	case TriggerTypeDevice:
		if t.Device == nil {
			return errors.Parameter.AddMsg("设备类型触发器需要填写触发详情")
		}
		return t.Device.Validate()
	}
	return nil
}
