package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupCreateLogic {
	return &TaskGroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupCreateLogic) TaskGroupCreate(in *timedjob.TaskGroup) (*timedjob.Response, error) {
	err := relationDB.NewTaskGroupRepo(l.ctx).Insert(l.ctx, ToTaskGroupPo(in))
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
