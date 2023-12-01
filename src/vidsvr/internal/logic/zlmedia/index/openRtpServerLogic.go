package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenRtpServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenRtpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenRtpServerLogic {
	return &OpenRtpServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenRtpServerLogic) OpenRtpServer(req *types.IndexApiReq) (resp *types.IndexApiOpenRtpServerResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, OPENRTPSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiOpenRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
