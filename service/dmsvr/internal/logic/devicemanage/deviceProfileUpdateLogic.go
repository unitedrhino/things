package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceProfileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceProfileUpdateLogic {
	return &DeviceProfileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceProfileUpdateLogic) DeviceProfileUpdate(in *dm.DeviceProfile) (*dm.Empty, error) {
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
	old, err := relationDB.NewDeviceProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceProfileFilter{
		Code: in.Code,
		Device: devices.Core{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
	})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		old = &relationDB.DmDeviceProfile{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName, Code: in.Code}
	}
	old.Params = in.Params
	err = relationDB.NewDeviceProfileRepo(l.ctx).Update(l.ctx, old)
	return &dm.Empty{}, err
}
