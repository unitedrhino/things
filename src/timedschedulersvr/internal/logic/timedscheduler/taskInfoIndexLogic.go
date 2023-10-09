package timedschedulerlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timedschedulersvr/internal/logic"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"
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

// 获取任务信息列表
func (l *TaskInfoIndexLogic) TaskInfoIndex(in *timedscheduler.TaskInfoIndexReq) (*timedscheduler.TaskInfoIndexResp, error) {
	if l.svcCtx.SchedulerRun == false {
		return nil, errors.NotEnable
	}
	db := relationDB.NewTaskRepo(l.ctx)
	f := relationDB.TaskFilter{
		Group:   in.Group,
		Type:    in.Type,
		SubType: in.SubType,
		Name:    in.Name,
		Code:    in.Code,
	}
	resp, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*timedscheduler.TaskInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, ToTimedTaskPb(v))
	}
	return &timedscheduler.TaskInfoIndexResp{Total: total, List: info}, nil
}
