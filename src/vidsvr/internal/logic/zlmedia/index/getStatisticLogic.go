package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, GETSTATISTIC, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStatisticResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err

}
