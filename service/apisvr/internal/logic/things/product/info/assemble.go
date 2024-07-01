package info

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func productInfoToApi(ctx context.Context, v *dm.ProductInfo) *types.ProductInfo {
	if uc := ctxs.GetUserCtx(ctx); uc != nil && !uc.IsAdmin {
		v.Secret = ""        // 设备秘钥
		v.ProtocolConf = nil // 设备证书
	}
	return utils.Copy[types.ProductInfo](v)
}

func productInfoToRpc(req *types.ProductInfo) *dm.ProductInfo {
	return utils.Copy[dm.ProductInfo](req)
}
