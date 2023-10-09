package task

import (
	"context"

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
	// todo: add your logic here and delete this line

	return nil
}
