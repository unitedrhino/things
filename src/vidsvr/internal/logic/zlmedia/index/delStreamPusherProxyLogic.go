package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, DELSTREAMPUSHERPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiDelStreamProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
