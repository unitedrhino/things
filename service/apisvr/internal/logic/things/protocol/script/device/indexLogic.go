package device

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取协议插件列表
func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.ProtocolScriptDeviceIndexReq) (resp *types.ProtocolScriptDeviceIndexResp, err error) {
	ret, err := l.svcCtx.ProtocolM.ProtocolScriptDeviceIndex(l.ctx, utils.Copy[dm.ProtocolScriptDeviceIndexReq](req))
	if err != nil {
		return nil, err
	}
	resp = utils.Copy[types.ProtocolScriptDeviceIndexResp](ret)
	if req.WithDevice {
		for _, v := range resp.List {
			if v.DeviceName == "" {
				dev, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName})
				if err != nil {
					continue
				}
				v.Device = utils.Copy[types.DeviceInfo](dev)
				continue
			}
			pi, err := l.svcCtx.ProductCache.GetData(l.ctx, v.ProductID)
			if err != nil {
				continue
			}
			v.Product = utils.Copy[types.ProductInfo](pi)
		}
	}

	return resp, err
}
