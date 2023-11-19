package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupReadLogic {
	return &TaskGroupReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupReadLogic) TaskGroupRead(in *timedjob.CodeReq) (*timedjob.TaskGroup, error) {
	po, err := relationDB.NewTaskGroupRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TaskGroupFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return ToTaskGroupPb(po), nil
}
