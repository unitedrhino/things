package devicegrouplogic

import (
	"context"
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
func (l *GroupInfoReadLogic) GroupInfoRead(in *dm.GroupInfoReadReq) (*dm.GroupInfo, error) {
	dg, err := l.GiDB.FindOneByFilter(l.ctx, relationDB.GroupInfoFilter{GroupID: in.GroupID})
	if err != nil {
		return nil, err
	}
	return &dm.GroupInfo{
		GroupID:     dg.GroupID,
		GroupName:   dg.GroupName,
		ParentID:    dg.ParentID,
		Desc:        dg.Desc,
		CreatedTime: dg.CreatedTime.Unix(),
		Tags:        dg.Tags,
	}, nil
}
