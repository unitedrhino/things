package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGroupMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGroupMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGroupMultiCreateLogic {
	return &DeviceGroupMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 将设备加到多个分组中
func (l *DeviceGroupMultiCreateLogic) DeviceGroupMultiCreate(in *dm.DeviceGroupMultiSaveReq) (*dm.Empty, error) {
	if len(in.GroupIDs) == 0 {
		return &dm.Empty{}, nil
	}
	_, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
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
	err = relationDB.NewGroupDeviceRepo(l.ctx).MultiInsert(l.ctx, gds)
	if err != nil {
		return nil, err
	}
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		err := logic.UpdateDevGroupsTags(ctx, l.svcCtx, []devices.Core{{ProductID: in.ProductID, DeviceName: in.DeviceName}})
		if err != nil {
			logx.WithContext(ctx).Errorf("update device group tags error: %s", err.Error())
		}
	})
	return &dm.Empty{}, nil
}
