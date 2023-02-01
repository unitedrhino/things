package log

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperIndexLogic {
	return &OperIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperIndexLogic) OperIndex(req *types.SysLogOperIndexReq) (resp *types.SysLogOperIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
