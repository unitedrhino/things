package category

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *SchemaIndexLogic) SchemaIndex(req *types.WithID) (resp *types.ProductCategorySchemaIndexResp, err error) {
	ret, err := l.svcCtx.ProductM.ProductCategorySchemaIndex(l.ctx, &dm.ProductCategorySchemaIndexReq{ProductCategoryID: req.ID})
	return utils.Copy[types.ProductCategorySchemaIndexResp](ret), err
}
