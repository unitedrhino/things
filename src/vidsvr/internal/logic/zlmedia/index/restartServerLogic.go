package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type RestartServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestartServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestartServerLogic {
	return &RestartServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestartServerLogic) RestartServer(req *types.IndexApiReq) (resp *types.IndexApiRestartServerResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, RESTARTSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiRestartServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
