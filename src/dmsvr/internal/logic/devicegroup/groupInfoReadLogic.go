package devicegrouplogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoReadLogic {
	return &GroupInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分组信息详情
func (l *GroupInfoReadLogic) GroupInfoRead(in *dm.GroupInfoReadReq) (*dm.GroupInfo, error) {
	dg, err := l.svcCtx.GroupInfo.FindOne(l.ctx, in.GroupID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, err
	}

	var tags map[string]string
	if dg.Tags.String != "" {
		_ = json.Unmarshal([]byte(dg.Tags.String), &tags)
	}
	return &dm.GroupInfo{
		GroupID:     dg.GroupID,
		GroupName:   dg.GroupName,
		ParentID:    dg.ParentID,
		Desc:        dg.Desc,
		CreatedTime: dg.CreatedTime.Unix(),
		Tags:        tags,
	}, nil
}
