package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStreamProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStreamProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStreamProxyLogic {
	return &AddStreamProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStreamProxyLogic) AddStreamProxy(req *types.IndexApiReq) (resp *types.IndexApiAddStreamProxyResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, ADDSTREAMPROXY, req.VidmgrID)
	dataRecv := new(types.IndexApiAddStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
