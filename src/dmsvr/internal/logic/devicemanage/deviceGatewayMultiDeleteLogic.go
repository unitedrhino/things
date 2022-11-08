package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGatewayMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayMultiDeleteLogic {
	return &DeviceGatewayMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除分组设备
func (l *DeviceGatewayMultiDeleteLogic) DeviceGatewayMultiDelete(in *dm.DeviceGatewayMultiDeleteReq) (*dm.Response, error) {
	err := l.svcCtx.Gateway.DeleteList(l.ctx, &device.Core{
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}, ToDeviceCoreDos(in.List))
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
