package appDeviceEvent

import (
	"context"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type AppDeviceHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewAppDeviceHandle(ctx context.Context, svcCtx *svc.ServiceContext) *AppDeviceHandle {
	return &AppDeviceHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (a *AppDeviceHandle) DeviceEventReport(in *application.EventReport) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	return nil
}

func (a *AppDeviceHandle) DevicePropertyReport(in *application.PropertyReport) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	infos, err := a.svcCtx.SceneDeviceRepo.GetInfos(a.ctx, in.Device, scene.DeviceOperationOperatorReportProperty, in.Identifier)
	if err != nil {
		a.Errorf("%s.GetInfos err:%v", err)
		return err
	}
	var exeInfos scene.Infos
	for _, info := range infos {
		if !info.Trigger.Device.IsTriggerWithProperty(in) {
			a.Infof("%s req=%v IsTriggerWithProperty not commit scene id:%v", utils.FuncName(), in, info.ID)
			continue
		}
		if len(info.When) != 0 {
			if !info.When.IsHit(a.ctx, scene.TermRepo{
				DeviceMsg:  a.svcCtx.DeviceMsg,
				SchemaRepo: a.svcCtx.SchemaRepo,
			}) {
				a.Infof("%s req=%v when not commit scene id:%v", utils.FuncName(), in, info.ID)
				continue
			}
		}
		exeInfos = append(exeInfos, info)
	}
	a.executeActions(exeInfos)
	return nil
}

func (a *AppDeviceHandle) DeviceStatusConnected(in *application.ConnectMsg) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	infos, err := a.svcCtx.SceneDeviceRepo.GetInfos(a.ctx, in.Device, scene.DeviceOperationOperatorConnected, "")
	if err != nil {
		a.Errorf("%s.GetInfos err:%v", err)
		return err
	}
	var exeInfos scene.Infos
	for _, info := range infos {
		if !info.Trigger.Device.IsTriggerWithConn(in.Device, scene.DeviceOperationOperatorConnected) {
			a.Infof("%s req=%v IsTriggerWithConn not commit scene id:%v", utils.FuncName(), in, info.ID)
			continue
		}
		if len(info.When) != 0 {
			if !info.When.IsHit(a.ctx, scene.TermRepo{
				DeviceMsg:  a.svcCtx.DeviceMsg,
				SchemaRepo: a.svcCtx.SchemaRepo,
			}) {
				a.Infof("%s req=%v when not commit scene id:%v", utils.FuncName(), in, info.ID)
				continue
			}
		}
		exeInfos = append(exeInfos, info)
	}
	a.executeActions(exeInfos)
	return nil
}

func (a *AppDeviceHandle) executeActions(exeInfos scene.Infos) {
	newCtx := utils.CopyContext(a.ctx)
	for _, info := range exeInfos {
		go func(ctx context.Context, info *scene.Info) (err error) {
			defer utils.Recover(ctx)
			startTime := time.Now().UnixMilli()
			defer logx.WithContext(ctx).Infof("%s.Execute end use:%vms sceneName:%v err:%v",
				utils.FuncName(), time.Now().UnixMilli()-startTime, info.Name, err)
			logx.WithContext(ctx).Infof("%s.Execute start sceneID:%v sceneName:%v", utils.FuncName(), info.ID, info.Name)
			err = info.Then.Execute(ctx, scene.ActionRepo{
				DeviceInteract: a.svcCtx.DeviceInteract,
				DeviceM:        a.svcCtx.DeviceM,
			})
			return err
		}(newCtx, info)
	}
}

func (a *AppDeviceHandle) DeviceStatusDisConnected(in *application.ConnectMsg) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	infos, err := a.svcCtx.SceneDeviceRepo.GetInfos(a.ctx, in.Device, scene.DeviceOperationOperatorDisConnected, "")
	if err != nil {
		a.Errorf("%s.GetInfos err:%v", err)
		return err
	}
	var exeInfos scene.Infos
	for _, info := range infos {
		if !info.Trigger.Device.IsTriggerWithConn(in.Device, scene.DeviceOperationOperatorDisConnected) {
			a.Infof("%s req=%v IsTriggerWithConn not commit scene id:%v", utils.FuncName(), in, info.ID)
			continue
		}
		if len(info.When) != 0 {
			if !info.When.IsHit(a.ctx, scene.TermRepo{
				DeviceMsg:  a.svcCtx.DeviceMsg,
				SchemaRepo: a.svcCtx.SchemaRepo,
			}) {
				a.Infof("%s req=%v when not commit scene id:%v", utils.FuncName(), in, info.ID)
				continue
			}
		}
		exeInfos = append(exeInfos, info)
	}
	a.executeActions(exeInfos)
	return nil
}
