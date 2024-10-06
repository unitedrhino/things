package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

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
	err := l.GdDB.MultiDelete(l.ctx, in.GroupID, logic.ToDeviceCores(in.List))
	if err != nil {
		return nil, err
	}
	relationDB.NewGroupInfoRepo(l.ctx).UpdateGroupDeviceCount(l.ctx, in.GroupID)
	return &dm.Empty{}, nil
}
