package schema

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiCreateLogic {
	return &MultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.ProductSchemaMultiCreateReq) error {
	dmReq := &dm.ProductSchemaMultiCreateReq{
		ProductID: req.ProductID,
		List:      ToSchemaInfosRpc(req.List),
	}
	_, err := l.svcCtx.ProductM.ProductSchemaMultiCreate(l.ctx, dmReq)
	return err
}
