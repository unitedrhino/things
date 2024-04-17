package category

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaMultiUpdateLogic {
	return &SchemaMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaMultiUpdateLogic) SchemaMultiUpdate(req *types.ProductCategorySchemaMultiUpdateReq) (err error) {
	_, err = l.svcCtx.ProductM.ProductCategorySchemaMultiUpdate(l.ctx, utils.Copy[dm.ProductCategorySchemaMultiUpdateReq](req))
	return err
}
