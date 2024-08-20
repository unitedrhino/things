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

type ProductCategorySchemaMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategorySchemaMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategorySchemaMultiCreateLogic {
	return &ProductCategorySchemaMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCategorySchemaMultiCreateLogic) ProductCategorySchemaMultiCreate(in *dm.ProductCategorySchemaMultiSaveReq) (*dm.Empty, error) {
	pcDB := relationDB.NewProductCategoryRepo(l.ctx)
	var productIDs []string
	{
		var productCategoryIDs = []int64{in.ProductCategoryID}
		var idPath string
		if in.ProductCategoryID != def.RootNode {
			pc, err := pcDB.FindOne(l.ctx, in.ProductCategoryID)
			if err != nil {
				return nil, err
			}
			idPath = pc.IDPath
		}
		pcs, err := pcDB.FindByFilter(l.ctx, relationDB.ProductCategoryFilter{IDPath: idPath}, nil)
		if err != nil {
			return nil, err
		}
		productCategoryIDs = append(productCategoryIDs, utils.ToSliceWithFunc(pcs, func(in *relationDB.DmProductCategory) int64 {
			return in.ID
		})...)
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
	oldIdentifierSet := utils.SliceToSet(oldIdentifiers)
	var newIdentifiers []string //把已经新增过的剔除
	for _, v := range in.Identifiers {
		if _, ok := oldIdentifierSet[v]; !ok { //原来没有的
			newIdentifiers = append(newIdentifiers, v)
		}
	}
	if len(newIdentifiers) == 0 {
		return &dm.Empty{}, nil
	}
	var insertDatas []*relationDB.DmProductCategorySchema
	for _, v := range newIdentifiers {
		insertDatas = append(insertDatas, &relationDB.DmProductCategorySchema{
			ProductCategoryID: in.ProductCategoryID,
			Identifier:        v,
		})
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewProductCategorySchemaRepo(l.ctx).MultiInsert(l.ctx, insertDatas)
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
				for _, v := range productIDs {
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
