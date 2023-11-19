package task

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupReadLogic {
	return &GroupReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupReadLogic) GroupRead(req *types.CodeReq) (resp *types.TimedTaskGroup, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskGroupRead(l.ctx, &timedjob.CodeReq{Code: req.Code})
	return ToGroupTypes(ret), err
}
