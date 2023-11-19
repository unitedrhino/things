package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupDeleteLogic {
	return &TaskGroupDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupDeleteLogic) TaskGroupDelete(in *timedjob.CodeReq) (*timedjob.Response, error) {
	err := relationDB.NewTaskGroupRepo(l.ctx).DeleteByFilter(l.ctx,
		relationDB.TaskGroupFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
