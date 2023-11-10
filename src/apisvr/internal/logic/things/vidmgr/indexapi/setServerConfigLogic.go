package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetServerConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetServerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetServerConfigLogic {
	return &SetServerConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetServerConfigLogic) SetServerConfig(req *types.IndexApiReq) (resp *types.IndexApiSetServerConfigResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, SETSERVERCONFIG, req.VidmgrID)
	dataRecv := new(types.IndexApiSetServerConfigResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
