package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimerReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimerReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimerReadLogic {
	return &DeviceTimerReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//func (l *DeviceTimerReadLogic) DeviceTimerRead(in *ud.WithID) (*ud.DeviceTimerInfo, error) {
//	po, err := relationDB.NewDeviceTimerInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
//	return ToDeviceTimerPb(po), err
//}
