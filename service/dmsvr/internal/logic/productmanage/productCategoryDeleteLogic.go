package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategoryDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategoryDeleteLogic {
	return &ProductCategoryDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除产品
func (l *ProductCategoryDeleteLogic) ProductCategoryDelete(in *dm.WithID) (*dm.Response, error) {
	err := relationDB.NewProductCategoryRepo(l.ctx).Delete(l.ctx, in.Id)
	return &dm.Response{}, err
}
