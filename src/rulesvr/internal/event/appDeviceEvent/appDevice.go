package appDeviceEvent

import (
	"context"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/utils"
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
	return nil
}

func (a *AppDeviceHandle) DeviceStatusConnected(in *application.ConnectMsg) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	return nil
}

func (a *AppDeviceHandle) DeviceStatusDisConnected(in *application.ConnectMsg) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	return nil
}
