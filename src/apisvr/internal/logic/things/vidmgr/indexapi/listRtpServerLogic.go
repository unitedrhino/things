package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRtpServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRtpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRtpServerLogic {
	return &ListRtpServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRtpServerLogic) ListRtpServer(req *types.IndexApiReq) (resp *types.IndexApiListRtpServerResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, LISTRTPSERVER, req.VidmgrID)
	dataRecv := new(types.IndexApiListRtpServerResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}
