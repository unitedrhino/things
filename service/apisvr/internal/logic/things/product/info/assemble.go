package info

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func productInfoToApi(v *dm.ProductInfo) *types.ProductInfo {
	return utils.Copy[types.ProductInfo](v)
}

func productInfoToRpc(req *types.ProductInfo) *dm.ProductInfo {
	if req == nil {
		return nil
	}
	return &dm.ProductInfo{
		ProductName:        req.ProductName,
		ProductID:          req.ProductID,
		AuthMode:           req.AuthMode,
		DeviceType:         req.DeviceType,
		CategoryID:         req.CategoryID,
		NetType:            req.NetType,
		Secret:             req.Secret, //动态注册产品秘钥 只读
		ProtocolCode:       req.ProtocolCode,
		AutoRegister:       req.AutoRegister,
		Desc:               utils.ToRpcNullString(req.Desc),
		Tags:               logic.ToTagsMap(req.Tags),
		ProductImg:         req.ProductImg,
		IsUpdateProductImg: req.IsUpdateProductImg,
	}
}
