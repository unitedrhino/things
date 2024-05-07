package interact

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayGetFoundSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayGetFoundSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayGetFoundSendLogic {
	return &GatewayGetFoundSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GatewayGetFoundSendLogic) GatewayGetFoundSend(req *types.GatewayGetFoundReq) error {
	_, err := l.svcCtx.DeviceInteract.GatewayGetFoundSend(l.ctx, utils.Copy[dm.GatewayGetFoundReq](req))

	return err
}
