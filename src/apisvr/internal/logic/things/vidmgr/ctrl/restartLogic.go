package ctrl

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestartLogic {
	return &RestartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestartLogic) Restart(req *types.CtrlApiReq) (resp *types.CtrlApiResp, err error) {
	// todo: add your logic here and delete this line

	return
}
