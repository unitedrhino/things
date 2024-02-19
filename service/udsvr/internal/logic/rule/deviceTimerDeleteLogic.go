package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimerDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimerDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimerDeleteLogic {
	return &DeviceTimerDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceTimerDeleteLogic) DeviceTimerDelete(in *ud.WithID) (*ud.Empty, error) {
	err := relationDB.NewDeviceTimerInfoRepo(l.ctx).Delete(l.ctx, in.Id)

	return &ud.Empty{}, err
}
