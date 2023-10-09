package task

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
