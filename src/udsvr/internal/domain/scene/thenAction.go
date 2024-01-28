// Package scene 执行动作
package scene

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// 操作执行器类型
type ActionExecutor string

const (
	ActionExecutorNotify ActionExecutor = "notify" //通知
	ActionExecutorDelay  ActionExecutor = "delay"  //延迟
	ActionExecutorDevice ActionExecutor = "device" //设备输出
	ActionExecutorAlarm  ActionExecutor = "alarm"  //告警
)

type Actions []*Action

type Action struct {
	Executor ActionExecutor `json:"executor"` //执行器类型 notify: 通知 delay:延迟  device:设备输出  alarm: 告警
	Delay    int64          `json:"delay"`    //秒数
	Alarm    *ActionAlarm   `json:"alarm"`
	Notify   *ActionNotify  `json:"notify"` //消息通知
	Device   *ActionDevice  `json:"device"`
}

func (a Actions) Validate() error {
	if a == nil {
		return nil
	}
	for _, v := range a {
		err := v.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Action) Validate() error {
	if a == nil {
		return nil
	}
	switch a.Executor {
	case ActionExecutorNotify:
		if a.Notify == nil {
			return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Executor))
		}
		return a.Notify.Validate()
	case ActionExecutorDelay:
		if a.Delay == 0 {
			return errors.Parameter.AddMsg("延时不能为0")
		}
	case ActionExecutorDevice:
		if a.Device == nil {
			return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Executor))
		}
		return a.Device.Validate()
	case ActionExecutorAlarm:
		if a.Alarm == nil {
			return errors.Parameter.AddMsg("对应的操作类型下没有进行配置:" + string(a.Executor))
		}
		return a.Alarm.Validate()
	default:
		return errors.Parameter.AddMsg("操作类型不支持:" + string(a.Executor))
	}
	return nil
}

// 执行操作
func (a *Action) Execute(ctx context.Context, repo ActionRepo) error {
	switch a.Executor {
	case ActionExecutorDelay:
		time.Sleep(time.Second * time.Duration(a.Delay))
	case ActionExecutorDevice:
		err := a.Device.Execute(ctx, repo)
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
