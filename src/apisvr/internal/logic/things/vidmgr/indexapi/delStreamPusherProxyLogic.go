package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelStreamPusherProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelStreamPusherProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelStreamPusherProxyLogic {
	return &DelStreamPusherProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelStreamPusherProxyLogic) DelStreamPusherProxy(req *types.IndexApiReq) (resp *types.IndexApiDelStreamProxyResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, DELSTREAMPUSHERPROXY, req.VidmgrID)
	dataRecv := new(types.IndexApiDelStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
