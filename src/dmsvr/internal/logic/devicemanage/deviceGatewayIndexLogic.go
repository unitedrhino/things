package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGatewayIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayIndexLogic {
	return &DeviceGatewayIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分组设备信息列表
func (l *DeviceGatewayIndexLogic) DeviceGatewayIndex(in *dm.DeviceGatewayIndexReq) (*dm.DeviceGatewayIndexResp, error) {
	f := mysql.GatewayDeviceFilter{
		Gateway: devices.Core{
			ProductID:  in.GatewayProductID,
			DeviceName: in.GatewayDeviceName,
		},
	}
	size, err := l.svcCtx.Gateway.CountByFilter(
		l.ctx, f)
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.Gateway.FindByFilter(
		l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info := make([]*dm.DeviceInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToDeviceInfo(v))
	}
	return &dm.DeviceGatewayIndexResp{List: info, Total: size}, nil
}
