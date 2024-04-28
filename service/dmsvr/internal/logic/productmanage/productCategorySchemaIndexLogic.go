package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategorySchemaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategorySchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategorySchemaIndexLogic {
	return &ProductCategorySchemaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品品类下的物模型列表,绑定的物模型会自动添加到该产品品类及子分类的产品中,并不支持删除
func (l *ProductCategorySchemaIndexLogic) ProductCategorySchemaIndex(in *dm.ProductCategorySchemaIndexReq) (*dm.ProductCategorySchemaIndexResp, error) {
	var ProductCategoryIDs = []int64{in.ProductCategoryID}
	if in.WithFather && in.ProductCategoryID != def.RootNode {
		pc, err := relationDB.NewProductCategoryRepo(l.ctx).FindOne(l.ctx, in.ProductCategoryID)
		if err != nil {
			return nil, err
		}
		ProductCategoryIDs = append(ProductCategoryIDs, utils.GetIDPath(pc.IDPath)...)
	}
	pos, err := relationDB.NewProductCategorySchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductCategorySchemaFilter{ProductCategoryIDs: ProductCategoryIDs}, nil)
	if err != nil {
		return nil, err
	}
	return &dm.ProductCategorySchemaIndexResp{
		Identifiers: utils.ToSliceWithFunc(pos, func(in *relationDB.DmProductCategorySchema) string {
			return in.Identifier
		})}, nil
}
