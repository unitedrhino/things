package category

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func productCategoryToApi(v *dm.ProductCategory) *types.ProductCategory {
	if v == nil {
		return nil
	}
	return &types.ProductCategory{
		ID:   v.Id,
		Name: v.Name,
		Desc: utils.ToNullString(v.Desc),
	}
}
func productCategoryToRpc(in *types.ProductCategory) *dm.ProductCategory {
	if in == nil {
		return nil
	}
	return &dm.ProductCategory{
		Id:   in.ID,
		Name: in.Name,
		Desc: utils.ToRpcNullString(in.Desc),
	}
}
