package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoUpdateLogic {
	return &GroupInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 更新分组
func (l *GroupInfoUpdateLogic) GroupInfoUpdate(in *dm.GroupInfoUpdateReq) (*dm.Response, error) {
	record, err := l.GiDB.FindOneByFilter(l.ctx, relationDB.GroupInfoFilter{GroupID: in.GroupID})
	if err != nil {
		return nil, err
	}

	err = l.GiDB.Update(l.ctx, &relationDB.DmGroupInfo{
		GroupID:   in.GroupID,
		ParentID:  record.ParentID,
		GroupName: in.GroupName,
		Desc:      in.Desc,
		Tags:      in.Tags,
	})
	if err != nil {
		return nil, errors.Parameter.AddMsg(err.Error())
	}

	return &dm.Response{}, nil
}
