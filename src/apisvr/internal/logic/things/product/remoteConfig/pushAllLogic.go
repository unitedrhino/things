package remoteConfig

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPushAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushAllLogic {
	return &PushAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushAllLogic) PushAll(req *types.ProductRemoteConfigPushAllReq) error {
	// todo: add your logic here and delete this line

	return nil
}
