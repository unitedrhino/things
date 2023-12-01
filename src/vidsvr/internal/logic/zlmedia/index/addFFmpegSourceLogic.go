package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFFmpegSourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFFmpegSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFFmpegSourceLogic {
	return &AddFFmpegSourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFFmpegSourceLogic) AddFFmpegSource(req *types.IndexApiReq) (resp *types.IndexApiAddFFmpegSourceResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, ADDFFMPEGSOURCE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiAddFFmpegSourceResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
