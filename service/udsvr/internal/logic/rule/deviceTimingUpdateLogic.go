package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

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
	old, err := relationDB.NewDeviceTimingInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.ExecAt = in.ExecAt
	old.Repeat = in.Repeat
	old.ActionType = in.ActionType
	old.DataID = in.DataID
	old.Value = in.Value
	old.Name = in.Name
	old.Status = in.Status
	err = DeviceTimingCheck(l.ctx, old)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewDeviceTimingInfoRepo(l.ctx).Update(l.ctx, old)
	return &ud.Empty{}, err
}
