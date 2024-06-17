package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimerUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimerUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimerUpdateLogic {
	return &DeviceTimerUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//func (l *DeviceTimerUpdateLogic) DeviceTimerUpdate(in *ud.DeviceTimerInfo) (*ud.Empty, error) {
//	old, err := relationDB.NewDeviceTimerInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
//	if err != nil {
//		return nil, err
//	}
//	old.ExecAt = in.ExecAt
//	old.ExecRepeat = in.ExecRepeat
//	old.ActionType = in.ActionType
//	old.DataID = in.DataID
//	old.Value = in.Value
//	old.Name = in.Name
//	old.Status = in.Status
//	old.LastRunTime = domain.GenLastRunTime(time.Now(), in.ExecAt)
//
//	err = DeviceTimerCheck(l.ctx, old)
//	if err != nil {
//		return nil, err
//	}
//	err = relationDB.NewDeviceTimerInfoRepo(l.ctx).Update(l.ctx, old)
//	return &ud.Empty{}, err
//}
