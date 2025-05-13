package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGroupMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGroupMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGroupMultiUpdateLogic {
	return &DeviceGroupMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新设备所在分组
func (l *DeviceGroupMultiUpdateLogic) DeviceGroupMultiUpdate(in *dm.DeviceGroupMultiSaveReq) (*dm.Empty, error) {
	if len(in.GroupIDs) == 0 {
		return &dm.Empty{}, nil
	}
	_, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		return nil, err
	}
	oldGs, err := relationDB.NewGroupInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.GroupInfoFilter{
		IDs: in.GroupIDs, Purpose: in.Purpose, HasDevice: &devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}}, nil)
	if err != nil {
		return nil, err
	}
	gs, err := relationDB.NewGroupInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.GroupInfoFilter{
		IDs: in.GroupIDs, Purpose: in.Purpose}, nil)
	if err != nil {
		return nil, err
	}
	var gds []*relationDB.DmGroupDevice
	for _, g := range gs {
		gds = append(gds, &relationDB.DmGroupDevice{
			GroupID:    g.ID,
			AreaID:     g.AreaID,
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
	}
	var oldGroupIDs []int64
	for _, v := range oldGs {
		oldGroupIDs = append(oldGroupIDs, v.ID)
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if len(oldGroupIDs) > 0 {
			err = relationDB.NewGroupDeviceRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.GroupDeviceFilter{
				GroupIDs:   oldGroupIDs,
				ProductID:  in.ProductID,
				DeviceName: in.DeviceName,
			})
			if err != nil {
				return err
			}
		}
		err = relationDB.NewGroupDeviceRepo(l.ctx).MultiInsert(l.ctx, gds)
		return err
	})
	if err != nil {
		return nil, err
	}
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := logic.UpdateDevGroupsTags(ctx, l.svcCtx, []devices.Core{{ProductID: in.ProductID, DeviceName: in.DeviceName}})
		if err != nil {
			logx.WithContext(ctx).Errorf("update device group tags error: %s", err.Error())
		}
	})
	return &dm.Empty{}, err
}
