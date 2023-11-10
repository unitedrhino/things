package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMediaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMediaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMediaListLogic {
	return &GetMediaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMediaListLogic) GetMediaList(req *types.IndexApiReq) (resp *types.IndexApiMediaListResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETMEDIALIST, req.VidmgrID)
	dataRecv := new(types.IndexApiMediaListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
