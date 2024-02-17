package rulelogic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &ud.DeviceTimingIndexResp{}, nil
}
