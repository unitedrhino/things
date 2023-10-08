package timedschedulerlogic

import (
	"context"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJobInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobInfoUpdateLogic {
	return &JobInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新任务
func (l *JobInfoUpdateLogic) JobInfoUpdate(in *timedscheduler.JobInfo) (*timedscheduler.Response, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.Response{}, nil
}
