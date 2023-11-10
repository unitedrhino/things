package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMediaPlayerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMediaPlayerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMediaPlayerListLogic {
	return &GetMediaPlayerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMediaPlayerListLogic) GetMediaPlayerList(req *types.IndexApiReq) (resp *types.IndexApiMediaPlayerListResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETMEDIAPLAYERLIST, req.VidmgrID)
	dataRecv := new(types.IndexApiMediaPlayerListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
