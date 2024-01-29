package productmanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/oss"
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
	if err != nil {
		return nil, err
	}
	if in.HeadImg != "" { //如果填了参数且不等于原来的,说明修改头像,需要处理
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProductManage, oss.SceneCategoryImg, fmt.Sprintf("%d/%s", po.ID, oss.GetFileNameWithPath(in.HeadImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.HeadImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.HeadImg = path
	}
	err = relationDB.NewProductCategoryRepo(l.ctx).Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	return &dm.WithID{Id: po.ID}, err
}
