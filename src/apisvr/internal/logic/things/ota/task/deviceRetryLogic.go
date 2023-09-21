package task

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceRetryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeviceRetryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceRetryLogic {
	return &DeviceRetryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceRetryLogic) DeviceRetry(req *types.OtaTaskDeviceRetryReq) error {
	_, err := l.svcCtx.DeviceMsg.OtaPromptIndex(l.ctx, &di.OtaPromptIndexReq{
		Id: req.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
