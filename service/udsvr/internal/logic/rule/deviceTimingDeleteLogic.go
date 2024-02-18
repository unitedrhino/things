package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimingDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimingDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimingDeleteLogic {
	return &DeviceTimingDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceTimingDeleteLogic) DeviceTimingDelete(in *ud.WithID) (*ud.Empty, error) {
	err := relationDB.NewDeviceTimingInfoRepo(l.ctx).Delete(l.ctx, in.Id)

	return &ud.Empty{}, err
}
