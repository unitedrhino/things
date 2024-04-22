// Package scene 执行动作
package scene

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// 操作执行器类型
type ActionType = string

const (
	ActionExecutorNotify ActionType = "notify" //通知 todo
	ActionExecutorDelay  ActionType = "delay"  //延迟
	ActionExecutorDevice ActionType = "device" //设备输出
	ActionExecutorAlarm  ActionType = "alarm"  //告警 todo
)

type Actions []*Action

type Action struct {
	ID     int64         `json:"id"`
	Order  int64         `json:"order"`
	Type   ActionType    `json:"type"`             //执行器类型 notify: 通知 delay:延迟  device:设备输出  alarm: 告警
	Delay  int64         `json:"delay,omitempty"`  //秒数
	Alarm  *ActionAlarm  `json:"alarm,omitempty"`  //todo
	Notify *ActionNotify `json:"notify,omitempty"` //消息通知 todo
	Device *ActionDevice `json:"device,omitempty"`
}

func (t *Action) GetFlowInfo() (ret *FlowInfo) {
	return &FlowInfo{
		Type:    "then",
		SubType: t.Type,
	}
}

func (a Actions) Validate(repo ValidateRepo) error {
	if a == nil {
		return nil
	}
	//var hasDevice bool
	for _, v := range a {
		//if v.Type == ActionExecutorDevice {
		//	hasDevice = true
		//}
		err := v.Validate(repo)
		if err != nil {
			return err
		}
	}
	//if !hasDevice {
	//	return errors.Parameter.AddMsg("执行必须有一个设备执行")
	//}
	return nil
}

func (a *Action) Validate(repo ValidateRepo) error {
	if a == nil {
		return nil
	}
	switch a.Type {
	//case ActionExecutorNotify:
	//	if a.Notify == nil {
	//		return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Type))
	//	}
	//	return a.Notify.Validate()
	case ActionExecutorDelay:
		if a.Delay == 0 {
			return errors.Parameter.AddMsg("延时不能为0")
		}
	case ActionExecutorDevice:
		if a.Device == nil {
			return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Type))
		}
		return a.Device.Validate(repo)
	case ActionExecutorNotify:
		if a.Notify == nil {
			return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Type))
		}
		return a.Notify.Validate(repo)
	//case ActionExecutorAlarm:
	//	if a.Alarm == nil {
	//		return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Type))
	//	}
	//	return a.Alarm.Validate()
	default:
		return errors.Parameter.AddMsg("操作类型不支持:" + string(a.Type))
	}
	return nil
}

// 执行操作
func (a *Action) Execute(ctx context.Context, repo ActionRepo) error {
	switch a.Type {
	case ActionExecutorDelay:
		time.Sleep(time.Second * time.Duration(a.Delay))
	case ActionExecutorDevice:
		err := a.Device.Execute(ctx, repo)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.Execute Action:%#v err:%v", utils.FuncName(), a, err)
			return err
		}
	case ActionExecutorNotify:
		err := a.Notify.Execute(ctx, repo)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.Execute Action:%#v err:%v", utils.FuncName(), a, err)
			return err
		}
	case ActionExecutorAlarm:
		err := a.Alarm.Execute(ctx, repo)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.Execute Action:%#v err:%v", utils.FuncName(), a, err)
			return err
		}
	}
	return nil
}

func (a Actions) Execute(ctx context.Context, repo ActionRepo) error {
	for _, v := range a {
		err := v.Execute(ctx, repo)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.Execute Action:%#v err:%v", utils.FuncName(), v, err)
			return err
		}
	}
	return nil
}
