package otaTask

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfinishedIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfinishedIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfinishedIndexLogic {
	return &UnfinishedIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnfinishedIndexLogic) UnfinishedIndex(req *types.OTAUnfinishedTaskByDeviceIndexReq) (resp *types.OTAUnfinishedTaskByDeviceIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
