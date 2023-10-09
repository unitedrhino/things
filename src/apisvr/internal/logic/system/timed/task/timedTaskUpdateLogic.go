package task

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TimedTaskUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTimedTaskUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TimedTaskUpdateLogic {
	return &TimedTaskUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TimedTaskUpdateLogic) TimedTaskUpdate(req *types.TimedTaskInfo) error {
	// todo: add your logic here and delete this line

	return nil
}
