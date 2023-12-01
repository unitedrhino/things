package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, CLOSERTPSERVER, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiCloseRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
