package info

import (
	"context"
	"github.com/i-Things/things/service/apisvr/internal/logic/things"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateLogic {
	return &MultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUpdateLogic) MultiUpdate(req *types.DeviceInfoMultiUpdateReq) error {
	_, err := l.svcCtx.DeviceM.DeviceInfoMultiUpdate(l.ctx, &dm.DeviceInfoMultiUpdateReq{
		Devices: things.ToDmDeviceCoresPb(req.Devices),
		AreaID:  req.AreaID,
	})
	return err
}
