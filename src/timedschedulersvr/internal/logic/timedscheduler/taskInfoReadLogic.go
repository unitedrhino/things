package timedschedulerlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

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

// 获取任务信息详情
func (l *TaskInfoReadLogic) TaskInfoRead(in *timedscheduler.TaskInfoReadReq) (*timedscheduler.TaskInfo, error) {
	if l.svcCtx.SchedulerRun == false {
		return nil, errors.NotEnable
	}
	db := relationDB.NewTaskRepo(l.ctx)
	po, err := db.FindOne(l.ctx, in.Id)
	if err == nil {
		return ToTimedTaskPb(po), nil
	}
	return &timedscheduler.TaskInfo{}, nil
}
