package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategoryReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategoryReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategoryReadLogic {
	return &ProductCategoryReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品信息详情
func (l *ProductCategoryReadLogic) ProductCategoryRead(in *dm.WithID) (*dm.ProductCategory, error) {
	po, err := relationDB.NewProductCategoryRepo(l.ctx).FindOne(l.ctx, in.Id)
	return ToProductCategoryRpc(l.ctx, po, l.svcCtx), err
}
