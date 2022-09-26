package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoIndexLogic {
	return &GroupInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分组信息列表
func (l *GroupInfoIndexLogic) GroupInfoIndex(in *dm.GroupInfoIndexReq) (*dm.GroupInfoIndexResp, error) {
	ros, total, err := l.svcCtx.GroupDB.Index(l.ctx, in)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	info := make([]*dm.GroupInfo, 0, len(ros))
	for _, ro := range ros {
		info = append(info, &dm.GroupInfo{
			GroupID:     ro.GroupID,
			GroupName:   ro.GroupName,
			ParentID:    ro.ParentID,
			Desc:        ro.Desc,
			CreatedTime: ro.CreatedTime,
			Tags:        in.Tags,
		})
	}

	rosAll, err := l.svcCtx.GroupDB.IndexAll(l.ctx, in)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	infoAll := make([]*dm.GroupInfo, 0, len(rosAll))
	for _, ro := range rosAll {
		infoAll = append(infoAll, &dm.GroupInfo{
			GroupID:     ro.GroupID,
			GroupName:   ro.GroupName,
			ParentID:    ro.ParentID,
			Desc:        ro.Desc,
			CreatedTime: ro.CreatedTime,
			Tags:        in.Tags,
		})
	}

	return &dm.GroupInfoIndexResp{List: info, Total: total, ListAll: infoAll}, nil
}
