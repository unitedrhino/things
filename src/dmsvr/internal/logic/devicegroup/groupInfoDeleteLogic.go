package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoDeleteLogic {
	return &GroupInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除分组
func (l *GroupInfoDeleteLogic) GroupInfoDelete(in *dm.GroupInfoDeleteReq) (*dm.Response, error) {
	//删除两表数据
	err := l.svcCtx.GroupDB.Delete(l.ctx, in.GroupID)
	if err != nil {
		return nil, errors.NotEmpty.AddDetail(err)
	}
	return &dm.Response{}, nil
}
