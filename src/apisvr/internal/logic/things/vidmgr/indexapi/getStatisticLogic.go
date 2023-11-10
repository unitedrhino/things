package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStatisticLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStatisticLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatisticLogic {
	return &GetStatisticLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStatisticLogic) GetStatistic(req *types.IndexApiReq) (resp *types.IndexApiStatisticResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETSTATISTIC, req.VidmgrID)
	dataRecv := new(types.IndexApiStatisticResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}
