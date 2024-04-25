package deviceinteractlogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayGetTopoSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayGetTopoSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayGetTopoSendLogic {
	return &GatewayGetTopoSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 实时获取网关拓扑关系
func (l *GatewayGetTopoSendLogic) GatewayGetTopoSend(in *dm.GatewayTopoReadSendReq) (*dm.GatewayTopoReadSendResp, error) {
	// todo: add your logic here and delete this line

	return &dm.GatewayTopoReadSendResp{}, nil
}
