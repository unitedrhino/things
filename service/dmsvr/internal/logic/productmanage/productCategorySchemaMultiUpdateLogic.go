package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategorySchemaMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategorySchemaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategorySchemaMultiUpdateLogic {
	return &ProductCategorySchemaMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCategorySchemaMultiUpdateLogic) ProductCategorySchemaMultiUpdate(in *dm.ProductCategorySchemaMultiUpdateReq) (*dm.Empty, error) {
	//需要将绑定的产品物模型进行调整为可选及必选
	err := relationDB.NewProductCategorySchemaRepo(l.ctx).MultiUpdate(l.ctx, in.ProductCategoryID, in.Identifiers)
	return &dm.Empty{}, err
}
