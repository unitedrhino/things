package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKickSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickSessionLogic {
	return &KickSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickSessionLogic) KickSession(req *types.IndexApiReq) (resp *types.IndexApiKickSessionResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, KICKSESSION, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiKickSessionResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
