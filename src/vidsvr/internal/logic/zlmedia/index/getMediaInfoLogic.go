package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, GETMEDIAINFO, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMediaInfoResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}
