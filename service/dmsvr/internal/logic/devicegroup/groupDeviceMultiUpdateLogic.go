package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeviceMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeviceMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeviceMultiUpdateLogic {
	return &GroupDeviceMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新分组设备
func (l *GroupDeviceMultiUpdateLogic) GroupDeviceMultiUpdate(in *dm.GroupDeviceMultiSaveReq) (*dm.Response, error) {
	t, err := relationDB.NewDeviceInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceFilter{Cores: logic.ToDeviceCores(in.List)})
	if err != nil {
		return nil, err
	}
	if int(t) != len(in.List) {
		return nil, errors.Duplicate.AddMsg("有被删除的设备请重试")
	}
	gi, err := relationDB.NewGroupInfoRepo(l.ctx).FindOne(l.ctx, in.GroupID)
	if err != nil {
		return nil, err
	}
	list := make([]*relationDB.DmGroupDevice, 0, len(in.List))
	for _, v := range in.List {
		list = append(list, &relationDB.DmGroupDevice{
			GroupID:    in.GroupID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			AreaID:     gi.AreaID,
		})
	}
	err = relationDB.NewGroupDeviceRepo(l.ctx).MultiUpdate(l.ctx, in.GroupID, list)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &dm.Response{}, nil
}
