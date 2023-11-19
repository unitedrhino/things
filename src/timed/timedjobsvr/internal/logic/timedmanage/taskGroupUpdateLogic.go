package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupUpdateLogic {
	return &TaskGroupUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupUpdateLogic) TaskGroupUpdate(in *timedjob.TaskGroup) (*timedjob.Response, error) {
	repo := relationDB.NewTaskGroupRepo(l.ctx)
	oldPo, err := repo.FindOneByFilter(l.ctx, relationDB.TaskGroupFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		oldPo.Name = in.Name
	}
	if in.Priority != 0 {
		oldPo.Priority = in.Priority
	}
	if in.Env != nil {
		oldPo.Env = in.Env
	}
	if in.Config != "" {
		oldPo.Config = in.Config
	}
	err = repo.Update(l.ctx, oldPo)
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
