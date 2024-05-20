package msg

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayCanBindIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayCanBindIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayCanBindIndexLogic {
	return &GatewayCanBindIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GatewayCanBindIndexLogic) GatewayCanBindIndex(req *types.GatewayCanBindIndexReq) (resp *types.GatewayCanBindIndexResp, err error) {
	ret, err := l.svcCtx.DeviceMsg.GatewayCanBindIndex(l.ctx, utils.Copy[dm.GatewayCanBindIndexReq](req))
	if err != nil {
		return nil, err
	}
	resp = &types.GatewayCanBindIndexResp{UpdatedTime: ret.UpdatedTime}
	//resp = utils.Copy[types.GatewayCanBindIndexResp](ret)
	for _, v := range ret.SubDevices {
		pi, err := l.svcCtx.ProductCache.GetData(l.ctx, v.ProductID)
		if err != nil {
			continue
		}
		resp.SubDevices = append(resp.SubDevices, &types.DeviceCore{
			ProductID:   v.ProductID,
			ProductName: pi.ProductName,
			DeviceName:  v.DeviceName,
		})
	}
	return resp, err
}
