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
	//查询是否存在子组，若存在 则无法直接删除，返回删除失败
	resp, err := l.svcCtx.GroupInfo.FindOneByParentID(l.ctx, in.GroupID)
	if resp != nil {
		return nil, errors.NotEmptyGroup.AddDetailf("the group have sun group can not delete.", in.GroupID)
	}

	//删除两表数据
	err = l.svcCtx.GroupDB.Delete(l.ctx, in.GroupID)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &dm.Response{}, nil
}
