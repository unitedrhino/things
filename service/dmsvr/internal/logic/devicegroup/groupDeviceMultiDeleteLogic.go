package devicegrouplogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeviceMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GdDB *relationDB.GroupDeviceRepo
}

func NewGroupDeviceMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeviceMultiDeleteLogic {
	return &GroupDeviceMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GdDB:   relationDB.NewGroupDeviceRepo(ctx),
	}
}

// 删除分组设备
func (l *GroupDeviceMultiDeleteLogic) GroupDeviceMultiDelete(in *dm.GroupDeviceMultiDeleteReq) (*dm.Empty, error) {
	if in.GroupID != 0 {
		err := l.GdDB.MultiDelete(l.ctx, in.GroupID, logic.ToDeviceCores(in.List))
		if err != nil {
			return nil, err
		}
		relationDB.NewGroupInfoRepo(l.ctx).UpdateGroupDeviceCount(l.ctx, in.GroupID)
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
			var devs []devices.Core
			for _, v := range in.List {
				devs = append(devs, devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName})
			}
			err = logic.UpdateDevGroupsTags(ctx, l.svcCtx, devs)
			if err != nil {
				logx.WithContext(ctx).Errorf("update device group tags error: %s", err.Error())
			}
		})
		return &dm.Empty{}, err
	}
	gs, err := relationDB.NewGroupInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.GroupInfoFilter{Purpose: in.Purpose, HasDevices: utils.CopySlice[devices.Core](in.List)}, nil)
	if err != nil {
		return nil, err
	}
	var groupIDs []int64
	for _, g := range gs {
		groupIDs = append(groupIDs, g.ID)
		err := l.GdDB.MultiDelete(l.ctx, g.ID, logic.ToDeviceCores(in.List))
		if err != nil {
			return nil, err
		}
		err = relationDB.NewGroupInfoRepo(l.ctx).UpdateGroupDeviceCount(l.ctx, g.ID)
		return &dm.Empty{}, err
	}
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		var devs []devices.Core
		for _, v := range in.List {
			devs = append(devs, devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName})
		}
		err = logic.UpdateDevGroupsTags(ctx, l.svcCtx, devs)
		if err != nil {
			logx.WithContext(ctx).Errorf("update device group tags error: %s", err.Error())
		}
	})

	return &dm.Empty{}, nil
}
