package gateway

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiCreateLogic {
	return &MultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.DeviceGateWayMultiCreateReq) error {
	m := make([]*dm.DeviceGatewayBindDevice, 0, len(req.List))
	for _, v := range req.List {
		m = append(m, &dm.DeviceGatewayBindDevice{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	_, err := l.svcCtx.DeviceM.DeviceGatewayMultiCreate(l.ctx,
		&dm.DeviceGatewayMultiCreateReq{
			Gateway: &dm.DeviceCore{
				ProductID:  req.GateWayProductID,
				DeviceName: req.GateWayDeviceName,
			},
			IsAuthSign: false,
			List:       m})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.MultiCreate req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
