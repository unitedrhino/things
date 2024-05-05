package devicemsglogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"

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
	ret, err := l.svcCtx.GatewayCanBind.GetDevices(l.ctx, devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	return utils.Copy[dm.GatewayCanBindIndexResp](ret), nil
}
