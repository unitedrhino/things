package timedmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"

	"github.com/i-Things/things/src/timed/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoIndexLogic {
	return &TaskInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoIndexLogic) TaskInfoIndex(in *timedjob.TaskInfoIndexReq) (*timedjob.TaskInfoIndexResp, error) {
	f := relationDB.TaskFilter{}
	repo := relationDB.NewTaskInfoRepo(l.ctx)
	list, err := repo.FindByFilter(l.ctx, f, ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &timedjob.TaskInfoIndexResp{Total: total, List: ToTaskInfoPbs(list)}, nil
}
