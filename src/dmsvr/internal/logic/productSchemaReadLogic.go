package logic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaReadLogic {
	return &ProductSchemaReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品物模型
func (l *ProductSchemaReadLogic) ProductSchemaRead(in *dm.ProductSchemaReadReq) (*dm.ProductSchema, error) {
	pt, err := l.svcCtx.SchemaRepo.GetSchemaInfo(l.ctx, in.ProductID)
	if err != nil {
		return nil, err
	}
	return ToProductSchema(pt), nil
}
