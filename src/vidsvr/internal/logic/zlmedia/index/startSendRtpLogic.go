package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type StartSendRtpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartSendRtpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartSendRtpLogic {
	return &StartSendRtpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartSendRtpLogic) StartSendRtp(req *types.IndexApiReq) (resp *types.IndexApiStartSendRtpResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, STARTSENDRTP, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStartSendRtpResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}
