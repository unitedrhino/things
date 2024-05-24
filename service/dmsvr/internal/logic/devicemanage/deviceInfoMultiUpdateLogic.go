package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoMultiUpdateLogic {
	return &DeviceInfoMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量更新设备状态
func (l *DeviceInfoMultiUpdateLogic) DeviceInfoMultiUpdate(in *dm.DeviceInfoMultiUpdateReq) (*dm.Empty, error) {
	if in.AreaID == def.RootNode {
		return nil, errors.Parameter.AddMsgf("设备不能在root节点的区域下")
	}
	for _, v := range in.Devices {
		//err := l.svcCtx.StatusRepo.ModifyDeviceArea(l.ctx, devices.Core{
		//	ProductID:  v.ProductID,
		//	DeviceName: v.DeviceName,
		//}, in.AreaID)
		//if err != nil {
		//	l.Error(err)
		//	return nil, errors.Database.AddDetail(err)
		//}
		//err = l.svcCtx.SendRepo.ModifyDeviceArea(l.ctx, devices.Core{
		//	ProductID:  v.ProductID,
		//	DeviceName: v.DeviceName,
		//}, in.AreaID)
		//if err != nil {
		//	l.Error(err)
		//	return nil, errors.Database.AddDetail(err)
		//}
		err := l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		}, nil)
		if err != nil {
			l.Error(err)
		}
	}
	err := relationDB.NewDeviceInfoRepo(l.ctx).MultiUpdate(l.ctx, logic.ToDeviceCores(in.Devices), &relationDB.DmDeviceInfo{AreaID: stores.AreaID(in.AreaID)})
	return &dm.Empty{}, err
}
