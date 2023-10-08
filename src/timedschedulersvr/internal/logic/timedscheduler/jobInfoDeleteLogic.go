package timedschedulerlogic

import (
	"context"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJobInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobInfoDeleteLogic {
	return &JobInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除任务
func (l *JobInfoDeleteLogic) JobInfoDelete(in *timedscheduler.JobInfoDeleteReq) (*timedscheduler.Response, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.Response{}, nil
}
