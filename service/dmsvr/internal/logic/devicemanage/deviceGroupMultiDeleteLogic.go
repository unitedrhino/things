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

type DeviceGroupMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceGroupMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGroupMultiDeleteLogic {
	return &DeviceGroupMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除设备所在分组
func (l *DeviceGroupMultiDeleteLogic) DeviceGroupMultiDelete(in *dm.DeviceGroupMultiSaveReq) (*dm.Empty, error) {
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
	var gids []int64
	for _, g := range gs {
		gids = append(gids, g.ID)
	}
	err = relationDB.NewGroupDeviceRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.GroupDeviceFilter{
		GroupIDs:   gids,
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
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
