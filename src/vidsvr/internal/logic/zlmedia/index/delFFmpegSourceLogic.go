package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, DELFFMPEGSOURCE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiDelFFmpegSourceResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
