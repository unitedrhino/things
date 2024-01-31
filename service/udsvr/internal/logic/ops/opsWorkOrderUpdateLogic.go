package opslogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpsWorkOrderUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpsWorkOrderUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpsWorkOrderUpdateLogic {
	return &OpsWorkOrderUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpsWorkOrderUpdateLogic) OpsWorkOrderUpdate(in *ud.OpsWorkOrder) (*ud.Empty, error) {
	old, err := relationDB.NewOpsWorkOrderRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	err = relationDB.NewOpsWorkOrderRepo(l.ctx).Update(l.ctx, old)
	return &ud.Empty{}, err
}
