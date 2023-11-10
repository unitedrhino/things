package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartSendRtpPassiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartSendRtpPassiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartSendRtpPassiveLogic {
	return &StartSendRtpPassiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartSendRtpPassiveLogic) StartSendRtpPassive(req *types.IndexApiReq) (resp *types.IndexApiStartSendRtpPassiveResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, STARTSENDRTPPASSIVE, req.VidmgrID)
	dataRecv := new(types.IndexApiStartSendRtpPassiveResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
