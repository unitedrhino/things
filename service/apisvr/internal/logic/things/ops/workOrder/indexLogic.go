package workOrder

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

func (l *IndexLogic) Index(req *types.OpsWorkOrderIndexReq) (resp *types.OpsWorkOrderIndexResp, err error) {
	ret, err := l.svcCtx.Ops.OpsWorkOrderIndex(l.ctx, &ud.OpsWorkOrderIndexReq{
		Page:   logic.ToUdPageRpc(req.Page),
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.OpsWorkOrder
	for _, v := range ret.List {
		list = append(list, &types.OpsWorkOrder{
			ID:          v.Id,
			AreaID:      v.AreaID,
			RaiseUserID: v.RaiseUserID,
			IssueDesc:   v.IssueDesc,
			Number:      v.Number,
			Type:        v.Type,
			Params:      v.Params,
			Status:      v.Status,
			CreatedTime: v.CreatedTime,
		})
	}
	return &types.OpsWorkOrderIndexResp{
		Total: ret.Total,
		List:  list,
	}, nil
}
