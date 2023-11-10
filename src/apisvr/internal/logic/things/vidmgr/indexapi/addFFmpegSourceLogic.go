package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	data, err := proxyMediaServer(l.ctx, l.svcCtx, ADDFFMPEGSOURCE, req.VidmgrID)
	dataRecv := new(types.IndexApiAddFFmpegSourceResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
