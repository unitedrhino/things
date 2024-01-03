package scene

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"time"
)

type Infos []*Info

type Info struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	CreatedTime time.Time `json:"createdTime"`
	Trigger     Trigger   `json:"trigger"` //多种触发方式
	When        When      `json:"when"`    //手动触发模式不生效
	Then        Then      `json:"then"`    //触发后执行的动作
	Status      int64     `json:"state"`   // 状态（1启用 2禁用）
}

func (i *Info) Validate() error {

	err := i.Trigger.Validate()
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
