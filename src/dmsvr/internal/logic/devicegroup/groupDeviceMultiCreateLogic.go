package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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
func (l *GroupDeviceMultiCreateLogic) GroupDeviceMultiCreate(in *dm.GroupDeviceMultiCreateReq) (*dm.Response, error) {
	//todo 需要检查设备是否存在

	list := make([]*relationDB.DmGroupDevice, 0, len(in.List))
	for _, v := range in.List {
		list = append(list, &relationDB.DmGroupDevice{
			GroupID:    in.GroupID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	err := l.GdDB.MultiInsert(l.ctx, list)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &dm.Response{}, nil
}
