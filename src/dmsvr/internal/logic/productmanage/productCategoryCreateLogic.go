package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategoryCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategoryCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategoryCreateLogic {
	return &ProductCategoryCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增产品
func (l *ProductCategoryCreateLogic) ProductCategoryCreate(in *dm.ProductCategory) (*dm.WithID, error) {
	po := relationDB.DmProductCategory{
		Name: in.Name,
		Desc: utils.ToEmptyString(in.Desc),
	}
	err := relationDB.NewProductCategoryRepo(l.ctx).Insert(l.ctx, &po)
	return &dm.WithID{Id: po.ID}, err
}
