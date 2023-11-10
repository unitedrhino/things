package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSnapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSnapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSnapLogic {
	return &GetSnapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSnapLogic) GetSnap(req *types.IndexApiReq) (resp *types.IndexApiSnapResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETSNAP, req.VidmgrID)
	dataRecv := new(types.IndexApiSnapResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
