package job

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DynamicCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDynamicCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DynamicCreateLogic {
	return &DynamicCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DynamicCreateLogic) DynamicCreate(req *types.DynamicUpgradeJobReq) (resp *types.UpgradeJobResp, err error) {
	// todo: add your logic here and delete this line

	return
}
