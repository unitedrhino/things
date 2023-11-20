package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiListLogic {
	return &GetApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApiListLogic) GetApiList(req *types.IndexApiReq) (resp *types.IndexApiListResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETAPILIST, req.VidmgrID)
	dataRecv := new(types.IndexApiListResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
