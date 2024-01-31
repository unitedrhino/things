package scene

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"time"
)

type TriggerType string

const (
	TriggerTypeDevice TriggerType = "device"
	TriggerTypeTimer  TriggerType = "timer"
	TriggerTypeManual TriggerType = "manual"
)

type Infos []*Info

type Info struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Desc        string      `json:"desc"`
	CreatedTime time.Time   `json:"createdTime"`
	TriggerType TriggerType `json:"triggerType"` //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Trigger     Trigger     `json:"trigger"`     //多种触发方式
	When        Terms       `json:"when"`
	Then        Actions     `json:"then"`
	Status      int64       `json:"state"` // 状态（1启用 2禁用）
}

func (i *Info) Validate() error {
	if !utils.SliceIn(i.TriggerType, TriggerTypeDevice, TriggerTypeTimer, TriggerTypeManual) {
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
	if i.Status == 0 {
		i.Status = def.Enable
	}
	return nil
}

type FindWithTriggerDto struct {
	devices.Core
	//Operator OperationSchema //触发类型  online:上线 offline:下线 reportProperty:属性上报 reportEvent: 事件上报
}

//func FindWithDeviceTrigger(ctx context.Context, svcCtx svc.ServiceContext, dot FindWithTriggerDto) []*Info {
//	return nil
//}
