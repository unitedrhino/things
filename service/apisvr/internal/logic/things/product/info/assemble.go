package info

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func productInfoToApi(v *dm.ProductInfo) *types.ProductInfo {
	return &types.ProductInfo{
		CreatedTime:  v.CreatedTime,              //创建时间 只读
		ProductID:    v.ProductID,                //产品id 只读
		ProductName:  v.ProductName,              //产品名称
		AuthMode:     v.AuthMode,                 //认证方式:0:账密认证,1:秘钥认证
		DeviceType:   v.DeviceType,               //设备类型:0:设备,1:网关,2:子设备
		CategoryID:   v.CategoryID,               //产品品类
		NetType:      v.NetType,                  //通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
		DataProto:    v.DataProto,                //数据协议:0:自定义,1:数据模板
		AutoRegister: v.AutoRegister,             //动态注册:0:关闭,1:打开,2:打开并自动创建设备
		Secret:       v.Secret,                   //动态注册产品秘钥 只读
		Desc:         utils.ToNullString(v.Desc), //描述
		Tags:         logic.ToTagsType(v.Tags),
		ProductImg:   v.ProductImg, //产品图片
	}
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
		DataProto:          req.DataProto,
		AutoRegister:       req.AutoRegister,
		Desc:               utils.ToRpcNullString(req.Desc),
		Tags:               logic.ToTagsMap(req.Tags),
		ProductImg:         req.ProductImg,
		IsUpdateProductImg: req.IsUpdateProductImg,
	}
}
