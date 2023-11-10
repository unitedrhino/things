package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStreamPusherProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStreamPusherProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStreamPusherProxyLogic {
	return &AddStreamPusherProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStreamPusherProxyLogic) AddStreamPusherProxy(req *types.IndexApiReq) (resp *types.IndexApiAddStreamPusherProxyResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, ADDSTREAMPUSHERPROXY, req.VidmgrID)
	dataRecv := new(types.IndexApiAddStreamPusherProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
