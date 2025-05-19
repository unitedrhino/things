package devicegrouplogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeviceMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GdDB *relationDB.GroupDeviceRepo
}

func NewGroupDeviceMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeviceMultiCreateLogic {
	return &GroupDeviceMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GdDB:   relationDB.NewGroupDeviceRepo(ctx),
	}
}

// 创建分组设备
func (l *GroupDeviceMultiCreateLogic) GroupDeviceMultiCreate(in *dm.GroupDeviceMultiSaveReq) (*dm.Empty, error) {
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
	t, err := relationDB.NewDeviceInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceFilter{Cores: logic.ToDeviceCores(in.List)})
	if err != nil {
		return nil, err
	}
	if int(t) != len(in.List) {
		return nil, errors.Parameter.AddMsg("有不存在的设备请重试")
	}
	gi, err := relationDB.NewGroupInfoRepo(l.ctx).FindOne(l.ctx, in.GroupID)
	if err != nil {
		return nil, err
	}
	gc := l.svcCtx.GroupConfig[gi.Purpose]

	list := make([]*relationDB.DmGroupDevice, 0, len(in.List))
	var devs []devices.Core
	for _, v := range in.List {
		list = append(list, &relationDB.DmGroupDevice{
			GroupID:    in.GroupID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			AreaID:     gi.AreaID,
		})
		devs = append(devs, devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName})
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if gc != nil && gc.UniqueDevice {
			err = relationDB.NewGroupDeviceRepo(tx).DeleteByFilter(l.ctx, relationDB.GroupDeviceFilter{
				Propose: gi.Purpose,
				Devs:    devs,
			})
			if err != nil {
				return err
			}
		}
		err = relationDB.NewGroupDeviceRepo(tx).MultiInsert(l.ctx, list)
		return err
	})
	if err != nil {
		return nil, err
	}
	relationDB.NewGroupInfoRepo(l.ctx).UpdateGroupDeviceCount(l.ctx, in.GroupID)

	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		ds, err := relationDB.NewGroupDeviceRepo(ctx).FindByFilter(ctx, relationDB.GroupDeviceFilter{GroupIDs: []int64{in.GroupID}, WithGroup: true}, nil)
		if err != nil {
			logx.WithContext(ctx).Errorf("dm.GroupDeviceMultiCreate err: %v", err)
			return
		}
		var devs []devices.Core
		for _, v := range ds {
			devs = append(devs, devices.Core{ProductID: v.ProductID, DeviceName: v.DeviceName})
		}
		err = logic.UpdateDevGroupsTags(ctx, l.svcCtx, devs)
		if err != nil {
			logx.WithContext(ctx).Errorf("update device group tags error: %s", err.Error())
		}
	})
	return &dm.Empty{}, nil
}
