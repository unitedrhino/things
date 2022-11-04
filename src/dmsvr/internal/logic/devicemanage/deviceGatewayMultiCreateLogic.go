package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGatewayMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayMultiCreateLogic {
	return &DeviceGatewayMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建分组设备
func (l *DeviceGatewayMultiCreateLogic) DeviceGatewayMultiCreate(in *dm.DeviceGatewayMultiCreateReq) (*dm.Response, error) {
	err := l.svcCtx.Gateway.CreateList(l.ctx, &device.Core{
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}, ToDeviceCoreDos(in.List))
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
