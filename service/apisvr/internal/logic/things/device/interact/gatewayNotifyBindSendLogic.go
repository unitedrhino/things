package interact

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayNotifyBindSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayNotifyBindSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayNotifyBindSendLogic {
	return &GatewayNotifyBindSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GatewayNotifyBindSendLogic) GatewayNotifyBindSend(req *types.GatewayNotifyBindSendReq) error {
	_, err := l.svcCtx.DeviceInteract.GatewayNotifyBindSend(l.ctx, utils.Copy[dm.GatewayNotifyBindSendReq](req))

	return err
}
