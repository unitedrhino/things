package productmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po := relationDB.DmProductCategory{
		ParentID: in.ParentID,
		Name:     in.Name,
		Desc:     utils.ToEmptyString(in.Desc),
	}

	err := relationDB.NewProductCategoryRepo(l.ctx).Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	if in.HeadImg != "" && in.IsUpdateHeadImg { //如果填了参数且不等于原来的,说明修改头像,需要处理
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProductManage, oss.SceneCategoryImg, fmt.Sprintf("%d/%s", po.ID, oss.GetFileNameWithPath(in.HeadImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.HeadImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.HeadImg = path
	}
	po.IDPath = cast.ToString(po.ID) + "-"
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if po.ParentID != 0 && po.ParentID != def.RootNode {
			parent, err := relationDB.NewProductCategoryRepo(tx).FindOne(l.ctx, in.ParentID)
			if err != nil {
				return err
			}
			po.IDPath = parent.IDPath + po.IDPath
			if parent.IsLeaf == def.True {
				parent.IsLeaf = def.False
				err = relationDB.NewProductCategoryRepo(tx).Update(l.ctx, parent)
				if err != nil {
					return err
				}
			}
		}
		err = relationDB.NewProductCategoryRepo(tx).Update(l.ctx, &po)
		return err
	})

	return &dm.WithID{Id: po.ID}, err
}
