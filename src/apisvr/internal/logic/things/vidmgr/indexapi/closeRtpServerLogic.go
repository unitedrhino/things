package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseRtpServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloseRtpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseRtpServerLogic {
	return &CloseRtpServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloseRtpServerLogic) CloseRtpServer(req *types.IndexApiReq) (resp *types.IndexApiCloseRtpServerResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, CLOSERTPSERVER, req.VidmgrID)
	dataRecv := new(types.IndexApiCloseRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
