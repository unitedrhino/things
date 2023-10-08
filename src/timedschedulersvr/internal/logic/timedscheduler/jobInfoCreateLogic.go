package timedschedulerlogic

import (
	"context"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJobInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobInfoCreateLogic {
	return &JobInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增任务
func (l *JobInfoCreateLogic) JobInfoCreate(in *timedscheduler.JobInfo) (*timedscheduler.Response, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.Response{}, nil
}
