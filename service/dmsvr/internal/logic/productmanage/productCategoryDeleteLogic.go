package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
func (l *ProductCategoryDeleteLogic) ProductCategoryDelete(in *dm.WithID) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po, err := relationDB.NewProductCategoryRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if po.ParentID != 0 {
			c, err := relationDB.NewProductCategoryRepo(tx).CountByFilter(l.ctx, relationDB.ProductCategoryFilter{ParentID: po.ParentID})
			if err != nil {
				return err
			}
			if c == 0 { //下面没有子节点了
				err = relationDB.NewProductCategoryRepo(tx).UpdateWithField(l.ctx,
					relationDB.ProductCategoryFilter{ID: po.ParentID}, map[string]any{"is_leaf": def.True})
				if err != nil {
					return err
				}
			}
		}
		err := relationDB.NewProductCategoryRepo(l.ctx).Delete(l.ctx, in.Id)
		return err
	})

	return &dm.Empty{}, err
}
