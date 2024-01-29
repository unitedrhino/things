package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/oss/common"
	"github.com/i-Things/things/src/dmsvr/internal/domain/productCustom"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToProductInfo(ctx context.Context, pi *relationDB.DmProductInfo, svcCtx *svc.ServiceContext) *dm.ProductInfo {

	if pi.DeviceType == def.Unknown {
		pi.DeviceType = def.DeviceTypeDevice
	}
	if pi.NetType == def.Unknown {
		pi.NetType = def.NetOther
	}
	if pi.DataProto == def.Unknown {
		pi.DataProto = def.DataProtoCustom
	}
	if pi.AuthMode == def.Unknown {
		pi.AuthMode = def.AuthModePwd
	}
	if pi.AutoRegister == def.Unknown {
		pi.AutoRegister = def.AutoRegClose
	}
	dpi := &dm.ProductInfo{
		ProductID:    pi.ProductID,                          //产品id
		ProductName:  pi.ProductName,                        //产品名
		AuthMode:     pi.AuthMode,                           //认证方式:0:账密认证,1:秘钥认证
		DeviceType:   pi.DeviceType,                         //设备类型:0:设备,1:网关,2:子设备
		CategoryID:   pi.CategoryID,                         //产品品类
		NetType:      pi.NetType,                            //通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
		DataProto:    pi.DataProto,                          //数据协议:0:自定义,1:数据模板
		AutoRegister: pi.AutoRegister,                       //动态注册:0:关闭,1:打开,2:打开并自动创建设备
		Secret:       pi.Secret,                             //动态注册产品秘钥 只读
		Desc:         &wrappers.StringValue{Value: pi.Desc}, //描述
		CreatedTime:  pi.CreatedTime.Unix(),                 //创建时间
		Tags:         pi.Tags,                               //产品tags
		ProductImg:   pi.ProductImg,
		//Model:     &wrappers.StringValue{Value: pi.Model},    //数据模板
	}
	if pi.ProductImg != "" {
		var err error
		dpi.ProductImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, pi.ProductImg, 24*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return dpi
}

func ToProductSchemaRpc(info *relationDB.DmProductSchema) *dm.ProductSchemaInfo {
	db := &dm.ProductSchemaInfo{
		ProductID:    info.ProductID,
		Tag:          info.Tag,
		Type:         info.Type,
		Identifier:   info.Identifier,
		ExtendConfig: info.ExtendConfig,
		Name:         utils.ToRpcNullString(&info.Name),
		Desc:         utils.ToRpcNullString(&info.Desc),
		Required:     info.Required,
		Affordance:   utils.ToRpcNullString(&info.Affordance),
	}
	return db
}

func ToProductSchemaPo(info *dm.ProductSchemaInfo) *relationDB.DmProductSchema {
	db := &relationDB.DmProductSchema{
		ProductID: info.ProductID,
		Tag:       info.Tag,
		DmSchemaCore: relationDB.DmSchemaCore{
			Type:         info.Type,
			Identifier:   info.Identifier,
			ExtendConfig: info.ExtendConfig,
			Name:         info.Name.GetValue(),
			Desc:         info.Desc.GetValue(),
			Required:     info.Required,
			Affordance:   info.Affordance.GetValue(),
		},
	}
	return db
}

func ToCustomTopicPb(info *productCustom.CustomTopic) *dm.CustomTopic {
	if info == nil {
		return nil
	}
	return &dm.CustomTopic{Topic: info.Topic, Direction: info.Direction}
}

func ToCustomTopicsPb(info []*productCustom.CustomTopic) (ret []*dm.CustomTopic) {
	for _, v := range info {
		ret = append(ret, ToCustomTopicPb(v))
	}
	return
}

func ToCustomTopicDo(info *dm.CustomTopic) *productCustom.CustomTopic {
	if info == nil {
		return nil
	}
	return &productCustom.CustomTopic{Topic: info.Topic, Direction: info.Direction}
}

func ToCustomTopicsDo(info []*dm.CustomTopic) (ret []*productCustom.CustomTopic) {
	for _, v := range info {
		ret = append(ret, ToCustomTopicDo(v))
	}
	return
}
func ToProductCategoryRpc(ctx context.Context, info *relationDB.DmProductCategory, svcCtx *svc.ServiceContext) *dm.ProductCategory {
	if info == nil {
		return nil
	}
	if info.HeadImg != "" {
		var err error
		info.HeadImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, info.HeadImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return &dm.ProductCategory{
		Id:      info.ID,
		Name:    info.Name,
		HeadImg: info.HeadImg,
		Desc:    utils.ToRpcNullString(info.Desc),
	}
}
