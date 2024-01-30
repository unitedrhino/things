package otaTask

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeLogic {
	return &UpgradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpgradeLogic) Upgrade(req *types.OTATaskReUpgradeReq) error {
	// todo: add your logic here and delete this line

	return nil
}
