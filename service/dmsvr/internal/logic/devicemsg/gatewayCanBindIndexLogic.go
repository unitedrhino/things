package devicemsglogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayCanBindIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayCanBindIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayCanBindIndexLogic {
	return &GatewayCanBindIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取网关可以绑定的子设备列表
func (l *GatewayCanBindIndexLogic) GatewayCanBindIndex(in *dm.GatewayCanBindIndexReq) (*dm.GatewayCanBindIndexResp, error) {
	gateway := devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}
	ret, err := l.svcCtx.GatewayCanBind.GetDevices(l.ctx, gateway)
	if err != nil {
		return nil, err
	}
	subDevices, err := relationDB.NewGatewayDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.GatewayDeviceFilter{Gateway: &gateway}, nil)
	if err != nil {
		return nil, err
	}
	var subMap map[devices.Core]struct{}
	for _, v := range subDevices {
		subMap[devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		}] = struct{}{}
	}
	var filterdDevs []*dm.DeviceCore
	for _, v := range ret.SubDevices {
		if _, ok := subMap[*v]; ok {
			continue
		}
		filterdDevs = append(filterdDevs, utils.Copy[dm.DeviceCore](v))
	}
	return &dm.GatewayCanBindIndexResp{
		SubDevices:  filterdDevs,
		UpdatedTime: ret.UpdatedTime,
	}, nil
}
