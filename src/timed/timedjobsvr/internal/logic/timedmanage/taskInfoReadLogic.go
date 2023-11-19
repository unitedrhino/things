package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoReadLogic {
	return &TaskInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoReadLogic) TaskInfoRead(in *timedjob.CodeReq) (*timedjob.TaskInfo, error) {
	po, err := relationDB.NewTaskInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TaskFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return ToTaskInfoPb(po), nil
}
