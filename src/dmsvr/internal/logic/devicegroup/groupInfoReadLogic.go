package devicegrouplogic

import (
	"context"
	devicemanagelogic "github.com/i-Things/things/src/dmsvr/internal/logic/devicemanage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoReadLogic {
	return &GroupInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 获取分组信息详情
func (l *GroupInfoReadLogic) GroupInfoRead(in *dm.WithID) (*dm.GroupInfo, error) {
	dg, err := l.GiDB.FindOne(l.ctx, in.Id)
	c, err := devicemanagelogic.NewDeviceInfoCountLogic(l.ctx, l.svcCtx).DeviceInfoCount(&dm.DeviceInfoCountReq{GroupIDs: []int64{dg.ID}})
	if err != nil {
		l.Error(err)
	}
	if err != nil {
		return nil, err
	}
	return ToGroupInfoPb(dg, c), nil
}
