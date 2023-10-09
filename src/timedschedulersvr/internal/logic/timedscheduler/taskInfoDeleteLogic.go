package timedschedulerlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoDeleteLogic {
	return &TaskInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除任务
func (l *TaskInfoDeleteLogic) TaskInfoDelete(in *timedscheduler.TaskInfoDeleteReq) (*timedscheduler.Response, error) {
	if l.svcCtx.SchedulerRun == false {
		return nil, errors.NotEnable
	}
	db := relationDB.NewTaskRepo(l.ctx)
	var po *relationDB.TimedTask
	var err error
	if in.Id != 0 {
		po, err = db.FindOne(l.ctx, in.Id)

	} else {
		po, err = db.FindOneByFilter(l.ctx, relationDB.TaskFilter{Code: in.Code, Group: in.Group})
	}
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.Scheduler.Unregister(po.EntryID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	err = db.Delete(l.ctx, po.ID)
	if err != nil {
		return nil, err
	}
	return &timedscheduler.Response{}, nil
}
