package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, ADDSTREAMPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiAddStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
