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

type TimedTaskDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTimedTaskDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TimedTaskDeleteLogic {
	return &TimedTaskDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TimedTaskDeleteLogic) TimedTaskDelete(req *types.TaskInfoDeleteReq) error {
	_, err := l.svcCtx.Timedscheduler.TaskInfoDelete(l.ctx, &timedscheduler.TaskInfoDeleteReq{
		Id: req.ID,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.TaskInfoDelete req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
