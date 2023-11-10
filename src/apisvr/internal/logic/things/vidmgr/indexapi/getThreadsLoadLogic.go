package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetThreadsLoadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetThreadsLoadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetThreadsLoadLogic {
	return &GetThreadsLoadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetThreadsLoadLogic) GetThreadsLoad(req *types.IndexApiReq) (resp *types.IndexApiThreadLoadResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETTHREADSLOAD, req.VidmgrID)
	dataRecv := new(types.IndexApiThreadLoadResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
