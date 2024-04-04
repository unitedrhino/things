package logic

import (
	"context"
	"gitee.com/i-Things/share/oss/common"
	"github.com/i-Things/things/service/dmsvr/internal/domain/productCustom"
	"github.com/i-Things/things/service/dmsvr/internal/domain/protocol"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"

	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToProductInfo(ctx context.Context, svcCtx *svc.ServiceContext, pi *relationDB.DmProductInfo) *dm.ProductInfo {
	if pi == nil {
		return nil
	}
	if pi.DeviceType == def.Unknown {
		pi.DeviceType = def.DeviceTypeDevice
	}
	if pi.NetType == def.Unknown {
		pi.NetType = def.NetOther
	}
	if pi.ProtocolCode == "" {
		pi.ProtocolCode = protocol.CodeIThings
	}
	if pi.AuthMode == def.Unknown {
		pi.AuthMode = def.AuthModePwd
	}
	if pi.AutoRegister == def.Unknown {
		pi.AutoRegister = def.AutoRegClose
	}
	dpi := utils.Copy[dm.ProductInfo](pi)
	dpi.Category = ToProductCategoryPb(ctx, svcCtx, pi.Category, nil)
	dpi.Protocol = ToProtocolInfoPb(pi.Protocol)
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
		ProductID:         info.ProductID,
		Tag:               info.Tag,
		Type:              info.Type,
		Identifier:        info.Identifier,
		ExtendConfig:      info.ExtendConfig,
		Name:              utils.ToRpcNullString(&info.Name),
		Desc:              utils.ToRpcNullString(&info.Desc),
		Required:          info.Required,
		IsCanSceneLinkage: info.IsCanSceneLinkage,
		IsShareAuthPerm:   info.IsShareAuthPerm,
		IsHistory:         info.IsHistory,
		Order:             info.Order,
		Affordance:        utils.ToRpcNullString(&info.Affordance),
	}
	return db
}

func ToProductSchemaPo(info *dm.ProductSchemaInfo) *relationDB.DmProductSchema {
	db := &relationDB.DmProductSchema{
		ProductID: info.ProductID,
		DmSchemaCore: relationDB.DmSchemaCore{
			Type:              info.Type,
			Identifier:        info.Identifier,
			ExtendConfig:      info.ExtendConfig,
			Name:              info.Name.GetValue(),
			Desc:              info.Desc.GetValue(),
			Required:          info.Required,
			IsCanSceneLinkage: info.IsCanSceneLinkage,
			IsShareAuthPerm:   info.IsShareAuthPerm,
			IsHistory:         info.IsHistory,
			Order:             info.Order,
			Affordance:        info.Affordance.GetValue(),
			Tag:               info.Tag,
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

func ToProductCategoryPb(ctx context.Context, svcCtx *svc.ServiceContext, info *relationDB.DmProductCategory, children []*relationDB.DmProductCategory) *dm.ProductCategory {
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
	ret := &dm.ProductCategory{
		Id:       info.ID,
		Name:     info.Name,
		HeadImg:  info.HeadImg,
		ParentID: info.ParentID,
		IdPath:   utils.GetIDPath(info.IDPath),
		Desc:     utils.ToRpcNullString(info.Desc),
	}
	if children != nil {
		var idMap = map[int64][]*dm.ProductCategory{}
		for _, v := range children {
			idMap[v.ParentID] = append(idMap[v.ParentID], ToProductCategoryPb(ctx, svcCtx, v, nil))
		}
		fillDictInfoChildren(ret, idMap)
	}
	return ret
}

func fillDictInfoChildren(node *dm.ProductCategory, nodeMap map[int64][]*dm.ProductCategory) {
	// 找到当前节点的子节点数组
	children := nodeMap[node.Id]
	for _, child := range children {
		fillDictInfoChildren(child, nodeMap)
	}
	node.Children = children
}
