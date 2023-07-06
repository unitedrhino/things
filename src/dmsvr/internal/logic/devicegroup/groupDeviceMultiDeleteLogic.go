package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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
func (l *GroupDeviceMultiDeleteLogic) GroupDeviceMultiDelete(in *dm.GroupDeviceMultiDeleteReq) (*dm.Response, error) {
	list := make([]*relationDB.DmGroupDevice, 0, len(in.List))
	for _, v := range in.List {
		if v == nil {
			continue
		}
		list = append(list, &relationDB.DmGroupDevice{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	err := l.GdDB.MultiDelete(l.ctx, in.GroupID, list)
	if err != nil {
		return nil, err
	}

	return &dm.Response{}, nil
}
