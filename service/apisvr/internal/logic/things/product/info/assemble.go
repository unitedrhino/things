package info

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func productInfoToApi(v *dm.ProductInfo) *types.ProductInfo {
	return utils.Copy[types.ProductInfo](v)
}

func productInfoToRpc(req *types.ProductInfo) *dm.ProductInfo {
	return utils.Copy[dm.ProductInfo](req)
}
