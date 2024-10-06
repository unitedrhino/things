package rulelogic

import (
	"context"
	"gitee.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimerIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimerIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimerIndexLogic {
	return &DeviceTimerIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//func (l *DeviceTimerIndexLogic) DeviceTimerIndex(in *ud.DeviceTimerIndexReq) (*ud.DeviceTimerIndexResp, error) {
//	f := relationDB.DeviceTimerInfoFilter{TriggerType: in.TriggerType, Msg: in.Msg}
//	total, err := relationDB.NewDeviceTimerInfoRepo(l.ctx).CountByFilter(l.ctx, f)
//	if err != nil {
//		return nil, err
//	}
//	list, err := relationDB.NewDeviceTimerInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
//	if err != nil {
//		return nil, err
//	}
//	return &ud.DeviceTimerIndexResp{List: ToDeviceTimersPb(list), Total: total}, nil
//}
