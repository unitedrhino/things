package category

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaMultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaMultiCreateLogic {
	return &SchemaMultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaMultiCreateLogic) SchemaMultiCreate(req *types.ProductCategorySchemaMultiSaveReq) error {
	_, err := l.svcCtx.ProductM.ProductCategorySchemaMultiCreate(l.ctx, utils.Copy[dm.ProductCategorySchemaMultiSaveReq](req))
	return err
}
