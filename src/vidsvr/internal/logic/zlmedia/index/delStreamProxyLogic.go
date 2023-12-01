package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelStreamProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelStreamProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelStreamProxyLogic {
	return &DelStreamProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelStreamProxyLogic) DelStreamProxy(req *types.IndexApiReq) (resp *types.IndexApiDelStreamProxyResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, DELSTREAMPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiDelStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
