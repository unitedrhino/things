package scene

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"time"
)

type Infos []*Info

type SceneType = string

const (
	//TriggerTypeDevice SceneType = "device"
	//TriggerTypeTimer  SceneType = "timer"
	SceneTypeManual SceneType = "manual" //手动触发 场景
	SceneTypeAuto   SceneType = "auto"   //自动化
)

type Info struct {
	ID          int64       `json:"id"`
	HeadImg     string      `json:"headImg"`  // 头像
	FlowPath    []*FlowInfo `json:"flowPath"` //执行路径
	Tag         string      `json:"tag"`
	Name        string      `json:"name"`
	Desc        string      `json:"desc"`
	CreatedTime time.Time   `json:"createdTime"`
	Type        SceneType   `json:"type"`
	If          If          `json:"if"`     //多种触发方式
	When        When        `json:"when"`   //手动触发模式不生效
	Then        Then        `json:"then"`   //触发后执行的动作
	Status      int64       `json:"status"` // 状态（1启用 2禁用）
}

func (i *Info) Validate(repo ValidateRepo) error {
	if !utils.SliceIn(i.Type, SceneTypeAuto, SceneTypeManual) {
		return errors.Parameter.AddMsg("场景类型不支持的类型:" + string(i.Type))
	}
	err := i.If.Validate(i.Type, repo)
	if err != nil {
		return err
	}
	err = i.When.Validate(repo)
	if err != nil {
		return err
	}
	err = i.Then.Validate(repo)
	if err != nil {
		return err
	}
	if i.Status == 0 {
		i.Status = def.Enable
	}
	i.FlowPath = i.Then.GetFlowPath()
	return nil
}

type FlowInfo struct {
	Type    string `json:"type"`    //流程类型 then
	SubType string `json:"subType"` //子类型 设备执行
	Info    string `json:"info"`    //设备执行类型为产品id
}

type FindWithTriggerDto struct {
	devices.Core
	//Type Schema //触发类型  online:上线 offline:下线 reportProperty:属性上报 reportEvent: 事件上报
}

//func FindWithDeviceTrigger(ctx context.Context, svcCtx svc.ServiceContext, dot FindWithTriggerDto) []*Info {
//	return nil
//}
