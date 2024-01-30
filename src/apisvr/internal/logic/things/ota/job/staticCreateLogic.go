package job

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StaticCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStaticCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StaticCreateLogic {
	return &StaticCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StaticCreateLogic) StaticCreate(req *types.StaticUpgradeJobReq) (resp *types.UpgradeJobResp, err error) {
	// todo: add your logic here and delete this line

	return
}
