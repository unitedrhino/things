package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeviceMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeviceMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeviceMultiCreateLogic {
	return &GroupDeviceMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建分组设备
func (l *GroupDeviceMultiCreateLogic) GroupDeviceMultiCreate(in *dm.GroupDeviceMultiCreateReq) (*dm.Response, error) {
	list := make([]*device.Core, len(in.List))
	for _, v := range in.List {
		list = append(list, &device.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	err := l.svcCtx.GroupDB.GroupDeviceCreate(l.ctx, in.GroupID, list)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &dm.Response{}, nil
}
