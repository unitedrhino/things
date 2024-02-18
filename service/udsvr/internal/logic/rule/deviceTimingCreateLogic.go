package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
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
	err := DeviceTimingCheck(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	relationDB.NewDeviceTimingInfoRepo(l.ctx).Insert(l.ctx, &po)

	return &ud.WithID{Id: po.ID}, nil
}

func DeviceTimingCheck(ctx context.Context, po *relationDB.UdDeviceTimingInfo) error {
	switch po.TriggerType {
	case deviceTiming.TriggerTimer:
		//todo 需要校验时间一样的
	case deviceTiming.TriggerDelay:
		count, err := relationDB.NewDeviceTimingInfoRepo(ctx).CountByFilter(ctx,
			relationDB.DeviceTimingInfoFilter{Devices: []*devices.Core{{po.ProductID, po.DeviceName}}, TriggerType: deviceTiming.TriggerDelay})
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.Duplicate.WithMsg("同时只能存在一个延时控制")
		}
	}
	return nil
}
