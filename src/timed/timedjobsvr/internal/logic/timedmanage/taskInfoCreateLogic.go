package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoCreateLogic {
	return &TaskInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoCreateLogic) TaskInfoCreate(in *timedjob.TaskInfo) (*timedjob.Response, error) {
	err := relationDB.NewTaskInfoRepo(l.ctx).Insert(l.ctx, ToTaskInfoPo(in))
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
