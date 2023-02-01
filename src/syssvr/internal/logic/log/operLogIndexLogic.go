package loglogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogIndexLogic {
	return &OperLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogIndexLogic) OperLogIndex(in *sys.OperLogIndexReq) (*sys.OperLogIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.OperLogIndexResp{}, nil
}
