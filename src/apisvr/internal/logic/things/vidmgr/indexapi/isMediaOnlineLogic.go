package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsMediaOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsMediaOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsMediaOnlineLogic {
	return &IsMediaOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsMediaOnlineLogic) IsMediaOnline(req *types.IndexApiReq) (resp *types.IndexApiIsMediaOnlineResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, ISMEDIAONLINE, req.VidmgrID)
	dataRecv := new(types.IndexApiIsMediaOnlineResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
