package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
