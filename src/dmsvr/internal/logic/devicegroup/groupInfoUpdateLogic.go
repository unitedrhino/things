package devicegrouplogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoUpdateLogic {
	return &GroupInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新分组
func (l *GroupInfoUpdateLogic) GroupInfoUpdate(in *dm.GroupInfoUpdateReq) (*dm.Response, error) {
	_, err := l.svcCtx.GroupInfo.FindOne(l.ctx, in.GroupID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind.AddDetailf("not find Group GroupID=%d",
				in.GroupID)
		}
		return nil, errors.Database.AddDetail(err)
	}

	var sqlTags sql.NullString
	if in.Tags != nil {
		tags, err := json.Marshal(in.Tags)
		if err == nil {
			sqlTags = sql.NullString{
				String: string(tags),
				Valid:  true,
			}
		}
	}

	err = l.svcCtx.GroupInfo.Update(l.ctx, &mysql.GroupInfo{
		GroupID:   in.GroupID,
		GroupName: in.GroupName,
		Desc:      in.Desc,
		Tags:      sqlTags.String,
	})

	return &dm.Response{}, nil
}
