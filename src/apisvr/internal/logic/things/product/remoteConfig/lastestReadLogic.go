package remoteConfig

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LastestReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLastestReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LastestReadLogic {
	return &LastestReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LastestReadLogic) LastestRead(req *types.ProductRemoteConfigLastestReadReq) (resp *types.ProductRemoteConfigLastestReadResp, err error) {
	// todo: add your logic here and delete this line

	return
}
