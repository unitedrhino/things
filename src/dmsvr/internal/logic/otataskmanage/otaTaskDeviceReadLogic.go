package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceReadLogic {
	return &OtaTaskDeviceReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设备升级状态详情
func (l *OtaTaskDeviceReadLogic) OtaTaskDeviceRead(in *dm.OtaTaskDeviceReadReq) (*dm.OtaTaskDeviceInfo, error) {
	var otDB = relationDB.NewOtaTaskDevicesRepo(l.ctx)
	otd, err := otDB.FindOne(l.ctx, in.ID)
	if err != nil {
		return nil, err
	}
	return ToOtaTaskDeviceInfo(otd), nil
}
