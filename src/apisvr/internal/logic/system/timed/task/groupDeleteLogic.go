package task

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteLogic {
	return &GroupDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupDeleteLogic) GroupDelete(req *types.CodeReq) error {
	l.Infof("req:%v", utils.Fmt(req))
	_, err := l.svcCtx.TimedJob.TaskGroupDelete(l.ctx, &timedjob.CodeReq{Code: req.Code})
	return err
}
