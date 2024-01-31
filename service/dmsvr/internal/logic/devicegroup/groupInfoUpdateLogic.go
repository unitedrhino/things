package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
func (l *GroupInfoUpdateLogic) GroupInfoUpdate(in *dm.GroupInfo) (*dm.Response, error) {
	record, err := l.GiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	record.Desc = in.Desc
	record.Name = in.Name
	record.Tags = in.Tags
	err = l.GiDB.Update(l.ctx, record)
	if err != nil {
		return nil, errors.Parameter.AddMsg(err.Error())
	}

	return &dm.Response{}, nil
}
