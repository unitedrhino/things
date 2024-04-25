package deviceinteractlogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayNotifyBindSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayNotifyBindSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayNotifyBindSendLogic {
	return &GatewayNotifyBindSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知网关绑定子设备
func (l *GatewayNotifyBindSendLogic) GatewayNotifyBindSend(in *dm.GatewayNotifyBindSendReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
