package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoDeleteLogic {
	return &GroupInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 删除分组
func (l *GroupInfoDeleteLogic) GroupInfoDelete(in *dm.WithID) (*dm.Response, error) {
	//删除两表数据
	err := l.GiDB.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
