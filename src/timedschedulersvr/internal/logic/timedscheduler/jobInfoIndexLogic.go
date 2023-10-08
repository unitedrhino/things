package timedschedulerlogic

import (
	"context"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJobInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobInfoIndexLogic {
	return &JobInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取任务信息列表
func (l *JobInfoIndexLogic) JobInfoIndex(in *timedscheduler.JobInfoIndexReq) (*timedscheduler.JobInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.JobInfoIndexResp{}, nil
}
