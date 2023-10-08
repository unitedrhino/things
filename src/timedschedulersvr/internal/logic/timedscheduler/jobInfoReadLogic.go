package timedschedulerlogic

import (
	"context"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJobInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobInfoReadLogic {
	return &JobInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取任务信息详情
func (l *JobInfoReadLogic) JobInfoRead(in *timedscheduler.JobInfoReadReq) (*timedscheduler.JobInfo, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.JobInfo{}, nil
}
