package opslogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsWorkOrderIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsWorkOrderIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsWorkOrderIndexLogic {
	return &OpsWorkOrderIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsWorkOrderIndexLogic) OpsWorkOrderIndex(in *ud.OpsWorkOrderIndexReq) (*ud.OpsWorkOrderIndexResp, error) {
	f := relationDB.OpsWorkOrderFilter{Status: in.Status}
	total, err := relationDB.NewOpsWorkOrderRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewOpsWorkOrderRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var retList []*ud.OpsWorkOrder
	for _, v := range list {
		retList = append(retList, &ud.OpsWorkOrder{
			Id:          v.ID,
			Number:      v.Number,
			RaiseUserID: v.RaiseUserID,
			Type:        v.Type,
			AreaID:      int64(v.AreaID),
			Params:      v.Params,
			IssueDesc:   v.IssueDesc,
			Status:      v.Status,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}

	return &ud.OpsWorkOrderIndexResp{List: retList, Total: total}, nil
}
