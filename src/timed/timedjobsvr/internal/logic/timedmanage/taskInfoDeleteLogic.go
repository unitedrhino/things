package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoDeleteLogic {
	return &TaskInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoDeleteLogic) TaskInfoDelete(in *timedjob.CodeReq) (*timedjob.Response, error) {
	err := relationDB.NewTaskInfoRepo(l.ctx).DeleteByFilter(l.ctx,
		relationDB.TaskFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
