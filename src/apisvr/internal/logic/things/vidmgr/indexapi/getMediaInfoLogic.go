package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMediaInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMediaInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMediaInfoLogic {
	return &GetMediaInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMediaInfoLogic) GetMediaInfo(req *types.IndexApiReq) (resp *types.IndexApiMediaInfoResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETMEDIAINFO, req.VidmgrID)
	dataRecv := new(types.IndexApiMediaInfoResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}
