package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllSessionLogic {
	return &GetAllSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllSessionLogic) GetAllSession(req *types.IndexApiReq) (resp *types.IndexApiAllSessionResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETALLSESSION, req.VidmgrID)
	dataRecv := new(types.IndexApiAllSessionResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
