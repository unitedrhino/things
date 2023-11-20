package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	data, err := proxyMediaServer(l.ctx, l.svcCtx, KICKSESSIONS, req.VidmgrID)
	dataRecv := new(types.IndexApiKickSessionsResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
