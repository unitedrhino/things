package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type StopSendRtpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopSendRtpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopSendRtpLogic {
	return &StopSendRtpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopSendRtpLogic) StopSendRtp(req *types.IndexApiReq) (resp *types.IndexApiStopSendRtpResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, STOPSENDRTP, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStopSendRtpResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
