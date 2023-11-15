package indexapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerConfigLogic {
	return &GetServerConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerConfigLogic) GetServerConfig(req *types.IndexApiReq) (resp *types.IndexApiServerConfigResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETSERVERCONFIG, req.VidmgrID)
	dataRecv := new(types.IndexApiServerConfigResp)
	json.Unmarshal(data, dataRecv)
	fmt.Println("GetServerConfig:", dataRecv)
	return dataRecv, err
}
