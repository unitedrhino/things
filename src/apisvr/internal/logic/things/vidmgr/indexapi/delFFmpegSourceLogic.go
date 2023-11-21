package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelFFmpegSourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelFFmpegSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFFmpegSourceLogic {
	return &DelFFmpegSourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelFFmpegSourceLogic) DelFFmpegSource(req *types.IndexApiReq) (resp *types.IndexApiDelFFmpegSourceResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, DELFFMPEGSOURCE, req.VidmgrID)
	dataRecv := new(types.IndexApiDelFFmpegSourceResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
