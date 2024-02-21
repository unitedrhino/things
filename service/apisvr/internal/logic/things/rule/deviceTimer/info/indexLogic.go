package info

import (
	"context"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DeviceTimerInfoIndexReq) (resp *types.DeviceTimerInfoIndexResp, err error) {
	ret, err := l.svcCtx.Rule.DeviceTimerIndex(l.ctx, &ud.DeviceTimerIndexReq{
		Page:        logic.ToUdPageRpc(req.Page),
		Status:      req.Status,
		TriggerType: req.TriggerType,
	})
	if err != nil {
		return nil, err
	}
	return &types.DeviceTimerInfoIndexResp{
		List:  ToInfosTypes(ret.List),
		Total: ret.Total,
	}, nil
}
