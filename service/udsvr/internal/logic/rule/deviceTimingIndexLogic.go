package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTimingIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTimingIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTimingIndexLogic {
	return &DeviceTimingIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceTimingIndexLogic) DeviceTimingIndex(in *ud.DeviceTimingIndexReq) (*ud.DeviceTimingIndexResp, error) {
	f := relationDB.DeviceTimingInfoFilter{TriggerType: in.TriggerType, Status: in.Status}
	total, err := relationDB.NewDeviceTimingInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewDeviceTimingInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &ud.DeviceTimingIndexResp{List: ToDeviceTimingsPb(list), Total: total}, nil
}
