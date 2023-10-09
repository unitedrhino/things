package task

import (
	"context"

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
	// todo: add your logic here and delete this line

	return nil
}
