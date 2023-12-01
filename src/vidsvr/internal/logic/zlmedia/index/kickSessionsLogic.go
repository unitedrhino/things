package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKickSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickSessionsLogic {
	return &KickSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickSessionsLogic) KickSessions(req *types.IndexApiReq) (resp *types.IndexApiKickSessionsResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, KICKSESSIONS, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiKickSessionsResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
