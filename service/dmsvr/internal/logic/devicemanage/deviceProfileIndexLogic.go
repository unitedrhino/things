package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceProfileIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceProfileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceProfileIndexLogic {
	return &DeviceProfileIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceProfileIndexLogic) DeviceProfileIndex(in *dm.DeviceProfileIndexReq) (*dm.DeviceProfileIndexResp, error) {
	pos, err := relationDB.NewDeviceProfileRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceProfileFilter{
		Codes: in.Codes,
		Device: devices.Core{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	var rets []*dm.DeviceProfile
	for _, po := range pos {
		ret := utils.Copy[dm.DeviceProfile](po)
		ret.Device = &dm.DeviceCore{
			ProductID:  po.ProductID,
			DeviceName: po.DeviceName,
		}
		rets = append(rets, ret)
	}
	return &dm.DeviceProfileIndexResp{Profiles: rets}, err
}
