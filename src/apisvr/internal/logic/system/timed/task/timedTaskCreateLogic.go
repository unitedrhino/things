package task

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TimedTaskCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTimedTaskCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TimedTaskCreateLogic {
	return &TimedTaskCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TimedTaskCreateLogic) TimedTaskCreate(req *types.TimedTaskInfo) error {
	_, err := l.svcCtx.Timedscheduler.TaskInfoCreate(l.ctx, ToTaskInfoPb(req))
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.TaskInfoCreate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
