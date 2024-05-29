package category

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ProductCategoryToApi(v *dm.ProductCategory) *types.ProductCategory {
	if v == nil {
		return nil
	}
	return &types.ProductCategory{
		ID:              v.Id,
		Name:            v.Name,
		Desc:            utils.ToNullString(v.Desc),
		HeadImg:         v.HeadImg,
		IDPath:          v.IdPath,
		ParentID:        v.ParentID,
		IsUpdateHeadImg: v.IsUpdateHeadImg,
		IsLeaf:          v.IsLeaf,
		Children:        productCategoriesToApi(v.Children),
	}
}

func productCategoriesToApi(in []*dm.ProductCategory) (ret []*types.ProductCategory) {
	for _, v := range in {
		ret = append(ret, ProductCategoryToApi(v))
	}
	return
}

func productCategoryToRpc(in *types.ProductCategory) *dm.ProductCategory {
	if in == nil {
		return nil
	}
	return &dm.ProductCategory{
		Id:              in.ID,
		ParentID:        in.ParentID,
		Name:            in.Name,
		Desc:            utils.ToRpcNullString(in.Desc),
		HeadImg:         in.HeadImg,
		IsUpdateHeadImg: in.IsUpdateHeadImg,
	}
}
