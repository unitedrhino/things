package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategorySchemaMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategorySchemaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategorySchemaMultiUpdateLogic {
	return &ProductCategorySchemaMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCategorySchemaMultiUpdateLogic) ProductCategorySchemaMultiUpdate(in *dm.ProductCategorySchemaMultiSaveReq) (*dm.Empty, error) {
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
	csDB := relationDB.NewCommonSchemaRepo(l.ctx)
	cs, err := csDB.FindByFilter(l.ctx, relationDB.CommonSchemaFilter{Identifiers: in.Identifiers}, nil)
	if err != nil {
		return nil, err
	}
	if len(cs) != len(in.Identifiers) {
		return nil, errors.Parameter.AddMsg("有物模型不存在")
	}
	pcsDB := relationDB.NewProductCategorySchemaRepo(l.ctx)
	olds, err := pcsDB.FindByFilter(l.ctx, relationDB.ProductCategorySchemaFilter{ProductCategoryID: in.ProductCategoryID}, nil)
	if err != nil {
		return nil, err
	}
	oldIdentifiers := utils.ToSliceWithFunc(olds, func(in *relationDB.DmProductCategorySchema) string {
		return in.Identifier
	})
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewProductSchemaRepo(tx).UpdateTag(l.ctx, productIDs, oldIdentifiers, schema.TagRequired, schema.TagOptional)
		if err != nil {
			return err
		}
		err = relationDB.NewProductCategorySchemaRepo(l.ctx).MultiUpdate(l.ctx, in.ProductCategoryID, in.Identifiers)
		if err != nil {
			return err
		}
		err = relationDB.NewProductSchemaRepo(tx).UpdateTag(l.ctx, productIDs, in.Identifiers, schema.TagOptional, schema.TagRequired)
		if err != nil {
			return err
		}
		return nil
	})
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
		adds := utils.GetAddSlice(oldIdentifiers, in.Identifiers)
		if len(adds) == 0 {
			return
		}
		addSet := utils.SliceToSet(adds)
		//物模型没有的需要增加
		for _, identifier := range cs {
			if _, ok := addSet[identifier.Identifier]; !ok {
				continue //非新增的直接跳过
			}
			identifier.ID = 0
			identifier.Tag = schema.TagRequired
			err = stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
				var findProducts = map[string]struct{}{} //数据库中有该物模型的
				psDB := relationDB.NewProductSchemaRepo(tx)
				//获取涉及到的产品物模型
				ps, err := psDB.FindProductIDByFilter(ctx, relationDB.ProductSchemaFilter{ProductIDs: productIDs, Identifiers: []string{identifier.Identifier}})
				if err != nil {
					return err
				}
				findProducts = utils.SliceToSet(ps)
				var schemas []*relationDB.DmProductSchema
				for _, v := range adds {
					if _, ok := findProducts[v]; ok {
						continue
					}
					//如果没有这个物模型需要新增
					schemas = append(schemas, &relationDB.DmProductSchema{
						ProductID:    v,
						DmSchemaCore: identifier.DmSchemaCore,
					})
				}
				return psDB.MultiInsert(ctx, schemas)
			})
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}
	})
	return &dm.Empty{}, err
}
