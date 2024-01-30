package otaTask

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeviceCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceCancelLogic {
	return &DeviceCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceCancelLogic) DeviceCancel(req *types.OTATaskByDeviceCancelReq) error {
	// todo: add your logic here and delete this line

	return nil
}
