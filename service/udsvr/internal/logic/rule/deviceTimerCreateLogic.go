package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/udsvr/internal/domain/deviceTimer"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimerCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimerCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimerCreateLogic {
	return &DeviceTimerCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//// 设备定时
//func (l *DeviceTimerCreateLogic) DeviceTimerCreate(in *ud.DeviceTimerInfo) (*ud.WithID, error) {
//	po := relationDB.UdDeviceTimerInfo{
//		ProductID:   in.Device.ProductID,
//		DeviceName:  in.Device.DeviceName,
//		TriggerType: in.TriggerType,
//		ExecAt:      in.ExecAt,
//		ExecRepeat:  in.ExecRepeat,
//		ActionType:  in.ActionType,
//		DataID:      in.DataID,
//		Value:       in.Value,
//		Name:        in.Name,
//		LastRunTime: domain.GenLastRunTime(time.Now(), in.ExecAt),
//		Msg:      in.Msg,
//	}
//	if po.Msg == 0 {
//		po.Msg = def.Enable
//	}
//	err := DeviceTimerCheck(l.ctx, &po)
//	if err != nil {
//		return nil, err
//	}
//	relationDB.NewDeviceTimerInfoRepo(l.ctx).Insert(l.ctx, &po)
//
//	return &ud.WithID{Id: po.ID}, nil
//}

func DeviceTimerCheck(ctx context.Context, po *relationDB.UdDeviceTimerInfo) error {
	switch po.TriggerType {
	case deviceTimer.TriggerTimer:
		//todo 需要校验时间一样的
	case deviceTimer.TriggerDelay:
		count, err := relationDB.NewDeviceTimerInfoRepo(ctx).CountByFilter(ctx,
			relationDB.DeviceTimerInfoFilter{Devices: []*devices.Core{{po.ProductID, po.DeviceName}}, TriggerType: deviceTimer.TriggerDelay})
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.Duplicate.WithMsg("同时只能存在一个延时控制")
		}
	}
	return nil
}
