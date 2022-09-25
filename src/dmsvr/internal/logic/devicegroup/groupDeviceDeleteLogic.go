package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeviceDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeviceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeviceDeleteLogic {
	return &GroupDeviceDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除分组设备
func (l *GroupDeviceDeleteLogic) GroupDeviceDelete(in *dm.GroupDeviceDeleteReq) (*dm.Response, error) {
	err := l.svcCtx.GroupDB.GroupDeviceDelete(l.ctx, in.GroupID, in.DeviceIndexList)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	return &dm.Response{}, nil
}
