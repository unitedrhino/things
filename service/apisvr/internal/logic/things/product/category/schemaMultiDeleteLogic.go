package category

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaMultiDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaMultiDeleteLogic {
	return &SchemaMultiDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaMultiDeleteLogic) SchemaMultiDelete(req *types.ProductCategorySchemaMultiSaveReq) error {
	_, err := l.svcCtx.ProductM.ProductCategorySchemaMultiDelete(l.ctx, utils.Copy[dm.ProductCategorySchemaMultiSaveReq](req))
	return err
}
