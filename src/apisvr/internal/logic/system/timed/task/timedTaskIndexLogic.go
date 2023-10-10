package task

import (
	"context"
	"github.com/i-Things/things/src/timedschedulersvr/client/timedscheduler"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TimedTaskIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTimedTaskIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TimedTaskIndexLogic {
	return &TimedTaskIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TimedTaskIndexLogic) TimedTaskIndex(req *types.TimedTaskIndexReq) (resp *types.TimedTaskIndexResp, err error) {
	var page timedscheduler.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.Timedscheduler.TaskInfoIndex(l.ctx, &timedscheduler.TaskInfoIndexReq{
		Page:    &page,
		Group:   req.Group,
		Type:    req.Type,
		SubType: req.SubType,
		Name:    req.Name,
		Code:    req.Code,
	})
	if err != nil {
		return nil, err
	}

	var total int64
	total = info.Total
	var tasks []*types.TimedTaskInfo
	tasks = make([]*types.TimedTaskInfo, 0, len(tasks))
	for _, i := range info.List {
		tasks = append(tasks, ToTaskInfoTypes(i))
	}
	return &types.TimedTaskIndexResp{List: tasks, Total: total}, nil
}
