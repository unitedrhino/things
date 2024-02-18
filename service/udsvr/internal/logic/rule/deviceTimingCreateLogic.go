package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/domain/deviceTiming"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

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
	po := relationDB.UdDeviceTimingInfo{
		ProductID:   in.Device.ProductID,
		DeviceName:  in.Device.DeviceName,
		TriggerType: in.TriggerType,
		ExecAt:      in.ExecAt,
		Repeat:      in.Repeat,
		ActionType:  in.ActionType,
		DataID:      in.DataID,
		Value:       in.Value,
		Name:        in.Name,
		Status:      in.Status,
	}
	switch po.TriggerType {
	case deviceTiming.TriggerTimer:
	case deviceTiming.TriggerDelay:
		relationDB.NewDeviceTimingInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceTimingInfoFilter{})
	}
	relationDB.NewDeviceTimingInfoRepo(l.ctx).Insert(l.ctx, &po)

	return &ud.WithID{}, nil
}
