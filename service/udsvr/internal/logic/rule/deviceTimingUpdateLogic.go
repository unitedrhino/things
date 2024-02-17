package rulelogic

import (
	"context"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimingUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimingUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimingUpdateLogic {
	return &DeviceTimingUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceTimingUpdateLogic) DeviceTimingUpdate(in *ud.DeviceTimingInfo) (*ud.Empty, error) {
	// todo: add your logic here and delete this line

	return &ud.Empty{}, nil
}
