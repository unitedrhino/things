package category

import (
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func productCategoryToApi(v *dm.ProductCategory) *types.ProductCategory {
	if v == nil {
		return nil
	}
	return &types.ProductCategory{
		ID:              v.Id,
		Name:            v.Name,
		Desc:            utils.ToNullString(v.Desc),
		HeadImg:         v.HeadImg,
		IsUpdateHeadImg: v.IsUpdateHeadImg,
	}
}
func productCategoryToRpc(in *types.ProductCategory) *dm.ProductCategory {
	if in == nil {
		return nil
	}
	return &dm.ProductCategory{
		Id:              in.ID,
		Name:            in.Name,
		Desc:            utils.ToRpcNullString(in.Desc),
		HeadImg:         in.HeadImg,
		IsUpdateHeadImg: in.IsUpdateHeadImg,
	}
}
