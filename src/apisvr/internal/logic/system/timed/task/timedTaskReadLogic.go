package task

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timedschedulersvr/pb/timedscheduler"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TimedTaskReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTimedTaskReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TimedTaskReadLogic {
	return &TimedTaskReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TimedTaskReadLogic) TimedTaskRead(req *types.TimedTaskReadReq) (resp *types.TimedTaskInfo, err error) {
	info, err := l.svcCtx.Timedscheduler.TaskInfoRead(l.ctx, &timedscheduler.TaskInfoReadReq{Id: req.ID})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.TaskInfoUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	return ToTaskInfoTypes(info), nil
}
