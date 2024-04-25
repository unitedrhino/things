package deviceinteractlogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayGetFoundSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayGetFoundSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayGetFoundSendLogic {
	return &GatewayGetFoundSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 实时获取网关拓扑关系
func (l *GatewayGetFoundSendLogic) GatewayGetFoundSend(in *dm.GatewayGetFoundReq) (*dm.GatewayGetFoundResp, error) {
	// todo: add your logic here and delete this line

	return &dm.GatewayGetFoundResp{}, nil
}
