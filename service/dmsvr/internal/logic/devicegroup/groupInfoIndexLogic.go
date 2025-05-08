package devicegrouplogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

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
		Names:       in.Names,
		ParentID:    in.ParentID,
		Tags:        in.Tags,
		Purpose:     in.Purpose,
		Purposes:    in.Purposes,
		WithProduct: true,
		AreaID:      in.AreaID,
		HasDevice:   utils.Copy[devices.Core](in.HasDevice),
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
		info = append(info, ToGroupInfoPb(l.ctx, l.svcCtx, ro))
	}
	return &dm.GroupInfoIndexResp{List: info, Total: total}, nil
}
