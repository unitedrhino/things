package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoCreateLogic {
	return &GroupInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建分组
func (l *GroupInfoCreateLogic) GroupInfoCreate(in *dm.GroupInfoCreateReq) (*dm.Response, error) {

	_, err := l.svcCtx.GroupInfo.Insert(l.ctx, &mysql.GroupInfo{
		GroupID:   l.svcCtx.GroupID.GetSnowflakeId(),
		ParentID:  in.ParentID,
		GroupName: in.GroupName,
		Desc:      in.Desc,
		Tags:      "{}",
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	return &dm.Response{}, nil
}
