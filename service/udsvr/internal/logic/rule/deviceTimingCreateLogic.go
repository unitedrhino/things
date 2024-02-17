package rulelogic

import (
	"context"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimingCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimingCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimingCreateLogic {
	return &DeviceTimingCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设备定时
func (l *DeviceTimingCreateLogic) DeviceTimingCreate(in *ud.DeviceTimingInfo) (*ud.WithID, error) {
	// todo: add your logic here and delete this line

	return &ud.WithID{}, nil
}
