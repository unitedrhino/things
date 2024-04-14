package devicemanagelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceBindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceBindLogic {
	return &DeviceBindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceBindLogic) DeviceBind(in *dm.DeviceBindReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
