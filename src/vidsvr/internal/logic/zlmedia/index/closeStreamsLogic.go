package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseStreamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloseStreamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseStreamsLogic {
	return &CloseStreamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloseStreamsLogic) CloseStreams(req *types.IndexApiReq) (resp *types.IndexApiCloseStreamsResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, CLOSESTREAMS, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiCloseStreamsResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
