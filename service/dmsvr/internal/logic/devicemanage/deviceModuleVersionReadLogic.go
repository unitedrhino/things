package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceModuleVersionReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceModuleVersionReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceModuleVersionReadLogic {
	return &DeviceModuleVersionReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceModuleVersionReadLogic) DeviceModuleVersionRead(in *dm.DeviceModuleVersionReadReq) (*dm.DeviceModuleVersion, error) {
	var (
		po  *relationDB.DmDeviceModuleVersion
		err error
	)

	if in.Id != 0 {
		po, err = relationDB.NewDeviceModuleVersionRepo(l.ctx).FindOne(l.ctx, in.Id)
	} else {
		po, err = relationDB.NewDeviceModuleVersionRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceModuleVersionFilter{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
			ModuleCode: in.ModuleCode,
		})
	}

	return utils.Copy[dm.DeviceModuleVersion](po), err
}
