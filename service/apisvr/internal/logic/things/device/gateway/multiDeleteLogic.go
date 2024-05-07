package gateway

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic/things"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiDeleteLogic {
	return &MultiDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MultiDeleteLogic) MultiDelete(req *types.DeviceGateWayMultiDeleteReq) error {
	_, err := l.svcCtx.DeviceM.DeviceGatewayMultiDelete(l.ctx,
		&dm.DeviceGatewayMultiSaveReq{
			Gateway: &dm.DeviceCore{
				ProductID:  req.GateWayProductID,
				DeviceName: req.GateWayDeviceName,
			},
			List: things.ToDmDeviceCoresPb(req.List)})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.MultiDelete MultiDelete req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
