package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloseStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseStreamLogic {
	return &CloseStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloseStreamLogic) CloseStream(req *types.IndexApiReq) (resp *types.IndexApiCloseStreamResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, CLOSESTREAM, req.VidmgrID)
	dataRecv := new(types.IndexApiCloseStreamResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
