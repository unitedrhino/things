package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceProfileReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceProfileReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceProfileReadLogic {
	return &DeviceProfileReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceProfileReadLogic) DeviceProfileRead(in *dm.DeviceProfileReadReq) (*dm.DeviceProfile, error) {
	po, err := relationDB.NewDeviceProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceProfileFilter{
		Code: in.Code,
		Device: devices.Core{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
	})
	if errors.Cmp(err, errors.NotFind) {
		return &dm.DeviceProfile{
			Device: in.Device,
			Code:   in.Code,
			Params: "",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	ret := utils.Copy[dm.DeviceProfile](po)
	ret.Device = &dm.DeviceCore{
		ProductID:  po.ProductID,
		DeviceName: po.DeviceName,
	}
	return ret, nil
}
