// Package scene 执行动作
package scene

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
)

type ActionDelayUnit string

const (
	ActionDelayUnitSeconds ActionDelayUnit = "seconds"
	ActionDelayUnitMinutes ActionDelayUnit = "minutes"
	ActionDelayUnitHours   ActionDelayUnit = "hours"
)

type ActionDelay struct {
	Time int64           `json:"time"` //延迟时间
	Unit ActionDelayUnit `json:"unit"` //时间单位 seconds:秒  minutes:分钟  hours:小时
}

func (a *ActionDelay) Validate() error {
	if a == nil {
		return nil
	}
	if !utils.SliceIn(a.Unit, ActionDelayUnitSeconds, ActionDelayUnitMinutes, ActionDelayUnitHours) {
		return errors.Parameter.AddMsg("操作延迟不支持的类型:" + string(a.Unit))
	}
	return nil
}
