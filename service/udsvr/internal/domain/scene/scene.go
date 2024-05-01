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

type DeviceMode = string

const (
	DeviceModeSingle DeviceMode = "single" //单设备
	DeviceModeMulti  DeviceMode = "multi"  //多设备
)

// 多设备的场景联动
type Info struct {
	ID          int64       `json:"id"`
	HeadImg     string      `json:"headImg"`               // 头像
	FlowPath    []*FlowInfo `json:"flowPath"`              //执行路径
	DeviceMode  DeviceMode  `json:"deviceMode"`            //设备模式: 1:单设备 2:多设备
	ProductID   string      `json:"productID,omitempty"`   //产品id
	DeviceName  string      `json:"deviceName,omitempty"`  //设备名
	DeviceAlias string      `json:"deviceAlias,omitempty"` //设备别名,只读
	Tag         string      `json:"tag"`
	Name        string      `json:"name"`
	Desc        string      `json:"desc"`
	CreatedTime time.Time   `json:"createdTime"`
	Type        SceneType   `json:"type"`
	If          If          `json:"if"`     //多种触发方式
	When        When        `json:"when"`   //手动触发模式不生效
	Then        Then        `json:"then"`   //触发后执行的动作
	Status      int64       `json:"status"` // 状态（1启用 2禁用）
	Body        string      `json:"body"`   //自定义字段
}

func (i *Info) Validate(repo ValidateRepo) error {
	if !utils.SliceIn(i.Type, SceneTypeAuto, SceneTypeManual) {
		return errors.Parameter.AddMsg("场景类型不支持的类型:" + string(i.Type))
	}
	if !utils.SliceIn(i.DeviceMode, DeviceModeSingle, DeviceModeMulti, "") {
		return errors.Parameter.AddMsg("场景设备模式不支持的类型:" + string(i.DeviceMode))
	}
	if i.DeviceMode == "" {
		i.DeviceMode = DeviceModeMulti
	}
	if i.DeviceMode == DeviceModeSingle && (i.ProductID == "" || i.DeviceName == "") {
		return errors.Parameter.AddMsg("单设备模式需要填写产品和设备")
	}
	if i.DeviceMode == DeviceModeSingle {
		i.DeviceAlias = GetDeviceAlias(repo.Ctx, repo.DeviceCache, i.ProductID, i.DeviceName)
	}
	repo.Info = i
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
