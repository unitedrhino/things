package task

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendDelayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendDelayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendDelayLogic {
	return &SendDelayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendDelayLogic) SendDelay(req *types.TimedTaskSendDelayReq) error {
	_, err := l.svcCtx.TimedJob.TaskSendDelay(l.ctx, ToSendDelayReqPb(req))
	return err
}
