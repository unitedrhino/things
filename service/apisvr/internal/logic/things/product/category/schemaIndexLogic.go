package category

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaIndexLogic {
	return &SchemaIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaIndexLogic) SchemaIndex(req *types.ProductCategorySchemaIndexReq) (resp *types.ProductCategorySchemaIndexResp, err error) {
	ret, err := l.svcCtx.ProductM.ProductCategorySchemaIndex(l.ctx, utils.Copy[dm.ProductCategorySchemaIndexReq](req))
	return utils.Copy[types.ProductCategorySchemaIndexResp](ret), err
}
