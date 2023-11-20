package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkThreadsLoadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkThreadsLoadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkThreadsLoadLogic {
	return &GetWorkThreadsLoadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkThreadsLoadLogic) GetWorkThreadsLoad(req *types.IndexApiReq) (resp *types.IndexApiWorkThreadLoadResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETWORKTHREADSLOAD, req.VidmgrID)
	dataRecv := new(types.IndexApiWorkThreadLoadResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
