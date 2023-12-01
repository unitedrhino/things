package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, GETMEDIALIST, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMediaListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
