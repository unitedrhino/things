package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, ADDSTREAMPUSHERPROXY, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiAddStreamPusherProxyResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
