package task

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceCancleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeviceCancleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceCancleLogic {
	return &DeviceCancleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceCancleLogic) DeviceCancle(req *types.OtaTaskDeviceCancleReq) error {
	_, err := l.svcCtx.OtaTaskM.OtaTaskDeviceCancle(l.ctx, &dm.OtaTaskDeviceCancleReq{
		ID: req.ID,
	})
	return err
}
