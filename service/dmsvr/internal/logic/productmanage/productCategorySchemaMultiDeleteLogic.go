package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategorySchemaMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategorySchemaMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategorySchemaMultiDeleteLogic {
	return &ProductCategorySchemaMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCategorySchemaMultiDeleteLogic) ProductCategorySchemaMultiDelete(in *dm.ProductCategorySchemaMultiSaveReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	pcDB := relationDB.NewProductCategoryRepo(l.ctx)
	var productIDs []string
	{
		var productCategoryIDs = []int64{in.ProductCategoryID}
		if in.ProductCategoryID != def.RootNode {
			pc, err := pcDB.FindOne(l.ctx, in.ProductCategoryID)
			if err != nil {
				return nil, err
			}
			productCategoryIDs = append(productCategoryIDs, utils.GetIDPath(pc.IDPath)...)
		}
		ps, err := relationDB.NewProductInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductFilter{CategoryIDs: productCategoryIDs}, nil)
		if err != nil {
			return nil, err
		}
		productIDs = utils.ToSliceWithFunc(ps, func(in *relationDB.DmProductInfo) string {
			return in.ProductID
		})
	}
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewProductSchemaRepo(tx).UpdateTag(l.ctx, productIDs, in.Identifiers, schema.TagRequired, schema.TagOptional)
		if err != nil {
			return err
		}
		err = relationDB.NewProductCategorySchemaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.ProductCategorySchemaFilter{ProductCategoryID: in.ProductCategoryID, Identifiers: in.Identifiers})
		if err != nil {
			return err
		}
		return nil
	})
	return &dm.Empty{}, err
}
