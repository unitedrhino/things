package timedschedulerlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoUpdateLogic {
	return &TaskInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新任务
func (l *TaskInfoUpdateLogic) TaskInfoUpdate(in *timedscheduler.TaskInfo) (*timedscheduler.Response, error) {
	if l.svcCtx.SchedulerRun == false {
		return nil, errors.NotEnable
	}
	db := relationDB.NewTaskRepo(l.ctx)
	po, err := db.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Error("RoleInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}
	if in.Group != "" {
		po.Group = in.Group
	}
	if in.Type != "" {
		po.Type = in.Type
	}
	if in.SubType != "" {
		po.SubType = in.SubType
	}
	if in.Name != "" {
		po.Name = in.Name
	}
	if in.Code != "" {
		po.Code = in.Code
	}
	if in.Params != "" {
		po.Params = in.Params
	}
	if in.CronExpression != "" {
		po.CronExpression = in.CronExpression
	}
	if in.Status != 0 {
		po.Status = in.Status
	}
	if in.Priority != "" {
		po.Priority = in.Priority
	}
	jb := PoToTaskDo(po)
	err = jb.Init()
	if err != nil {
		return nil, err
	}
	err = db.Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	l.svcCtx.Scheduler.Unregister(po.EntryID)
	if po.Status != 3 {
		task := jb.ToTask()
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
	}
	return &timedscheduler.Response{}, nil
}
