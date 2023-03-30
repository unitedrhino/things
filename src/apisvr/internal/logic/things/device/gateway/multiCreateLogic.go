package gateway

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.DeviceGateWayMultiCreateReq) error {
	m := make([]*dm.DeviceCore, 0, len(req.List))
	for _, v := range req.List {
		m = append(m, &dm.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	_, err := l.svcCtx.DeviceM.DeviceGatewayMultiCreate(l.ctx,
		&dm.DeviceGatewayMultiCreateReq{
			GatewayProductID:  req.GateWayProductID,
			GatewayDeviceName: req.GateWayDeviceName,
			List:              m})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.MultiCreate req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
