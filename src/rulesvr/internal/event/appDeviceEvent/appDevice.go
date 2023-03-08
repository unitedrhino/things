package appDeviceEvent

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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
	fmt.Println(infos)
	return nil
}

func (a *AppDeviceHandle) DeviceStatusConnected(in *application.ConnectMsg) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	infos, err := a.svcCtx.SceneDeviceRepo.GetInfos(a.ctx, in.Device, scene.DeviceOperationOperatorConnected, "")
	if err != nil {
		a.Errorf("%s.GetInfos err:%v", err)
		return err
	}
	for _, info := range infos {
		if len(info.When) != 0 {
			if !info.When.IsTrue(a.ctx, scene.TermRepo{
				DeviceInteract: a.svcCtx.DeviceInteract,
				DeviceMsg:      a.svcCtx.DeviceMsg,
				SchemaRepo:     a.svcCtx.SchemaRepo,
			}) {
				a.Infof("%s req=%v when not commit scene id:%v", utils.FuncName(), in, info.ID)
				return nil
			}
		}
	}
	fmt.Println(infos)
	return nil
}

func (a *AppDeviceHandle) DeviceStatusDisConnected(in *application.ConnectMsg) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	infos, err := a.svcCtx.SceneDeviceRepo.GetInfos(a.ctx, in.Device, scene.DeviceOperationOperatorDisConnected, "")
	if err != nil {
		a.Errorf("%s.GetInfos err:%v", err)
		return err
	}
	fmt.Println(infos)
	return nil
}
