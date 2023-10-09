package timedschedulerlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoCreateLogic {
	return &TaskInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增任务
func (l *TaskInfoCreateLogic) TaskInfoCreate(in *timedscheduler.TaskInfo) (*timedscheduler.Response, error) {
	if l.svcCtx.SchedulerRun == false {
		return nil, errors.NotEnable
	}
	db := relationDB.NewTaskRepo(l.ctx)
	jb := ToTaskDo(in)
	err := jb.Init()
	if err != nil {
		return nil, err
	}
	po := ToTimedTaskPo(in)

	err = db.Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	task := jb.ToTask()
	// every one minute exec
	entryID, err := l.svcCtx.Scheduler.Register(po.CronExpression, task)
	if err != nil {
		l.Errorf("Scheduler.Register  err:%+v , task:%+v", err, task)
		return nil, err
	}
	po.EntryID = entryID
	po.Status = relationDB.TaskStatusRun
	err = db.Update(l.ctx, po)
	if err != nil {
		l.Errorf("Scheduler.Update  err:%+v , task:%+v", err, task)
		l.svcCtx.Scheduler.Unregister(entryID)
	}
	return &timedscheduler.Response{}, nil
}
