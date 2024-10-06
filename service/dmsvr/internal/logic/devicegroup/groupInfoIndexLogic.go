package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoIndexLogic {
	return &GroupInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 获取分组信息列表
func (l *GroupInfoIndexLogic) GroupInfoIndex(in *dm.GroupInfoIndexReq) (*dm.GroupInfoIndexResp, error) {
	f := relationDB.GroupInfoFilter{
		Name:        in.Name,
		ParentID:    in.ParentID,
		Tags:        in.Tags,
		WithProduct: true,
		AreaID:      in.AreaID,
	}
	ros, err := l.GiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.GiDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*dm.GroupInfo, 0, len(ros))
	for _, ro := range ros {
		info = append(info, ToGroupInfoPb(ro))
	}
	return &dm.GroupInfoIndexResp{List: info, Total: total}, nil
}
