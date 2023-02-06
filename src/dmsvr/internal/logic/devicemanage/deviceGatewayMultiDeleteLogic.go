package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
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
	err := l.svcCtx.Gateway.DeleteList(l.ctx, &devices.Core{
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}, ToDeviceCoreDos(in.List))
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DataUpdate.DeviceGatewayUpdate(l.ctx, &events.GatewayUpdateInfo{
		GatewayProductID:  in.GatewayProductID,
		GatewayDeviceName: in.GatewayDeviceName,
		Status:            def.GatewayUnbind,
		Devices:           ToDeviceCoreEvents(in.List),
	})
	if err != nil {
		l.Errorf("%s.DeviceGatewayUpdate err=%+v", utils.FuncName(), err)
	}
	return &dm.Response{}, nil
}
