package scene

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"time"
)

type TriggerType string

const (
	TriggerTypeDevice TriggerType = "device"
	TriggerTypeTimer  TriggerType = "timer"
	TriggerTypeManual TriggerType = "manual"
)

type Info struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Desc        string      `json:"desc"`
	CreatedTime time.Time   `json:"createdTime"`
	TriggerType TriggerType `json:"triggerType"` //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Trigger     Trigger     `json:"trigger"`
	When        Terms       `json:"when"` //只有设备触发时才有用
	Then        Actions     `json:"then"`
	State       int64       `json:"state"` // 状态（1启用 2禁用）
}

func (i *Info) Validate() error {
	if !utils.SliceIn(i.TriggerType, TriggerTypeDevice) {
		return errors.Parameter.AddMsg("触发器不支持的类型:" + string(i.TriggerType))
	}
	err := i.Trigger.Validate(i.TriggerType)
	if err != nil {
		return err
	}
	err = i.When.Validate()
	if err != nil {
		return err
	}
	err = i.Then.Validate()
	if err != nil {
		return err
	}
	if i.State == 0 {
		i.State = def.Enable
	}
	return nil
}
