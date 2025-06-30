package logic

import (
	"context"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/productCustom"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
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
	if pi.OnlineHandle == def.Unknown {
		pi.OnlineHandle = def.True
	}
	var err error
	if len(pi.CustomUi) != 0 {
		for _, v := range pi.CustomUi {
			v.Path, err = svcCtx.OssClient.PublicBucket().GetUrl(v.Path, false)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.CustomUiGetUrl err:%v", utils.FuncName(), err)
			}
		}
	}

	dpi := utils.Copy[dm.ProductInfo](pi)
	dpi.Category = ToProductCategoryPb(ctx, svcCtx, pi.Category, nil)
	dpi.Protocol = ToProtocolInfoPb(pi.Protocol)
	if pi.ProductImg != "" {
		dpi.ProductImg, err = svcCtx.OssClient.PublicBucket().GetUrl(pi.ProductImg, false)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return dpi
}

func ToProductSchemaRpc(info *relationDB.DmSchemaInfo) *dm.ProductSchemaInfo {
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
		FuncGroup:         info.FuncGroup,
		ControlMode:       info.ControlMode,
		UserPerm:          info.UserPerm,
		IsHistory:         info.IsHistory,
		IsPassword:        info.IsPassword,
		Order:             info.Order,
		Affordance:        utils.ToRpcNullString(&info.Affordance),
	}
	return db
}

func ToProductSchemaPo(info *dm.ProductSchemaInfo) *relationDB.DmSchemaInfo {
	return utils.Copy[relationDB.DmSchemaInfo](info)
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
		Id:          info.ID,
		Name:        info.Name,
		HeadImg:     info.HeadImg,
		ParentID:    info.ParentID,
		IdPath:      utils.GetIDPath(info.IDPath),
		Desc:        utils.ToRpcNullString(info.Desc),
		IsLeaf:      info.IsLeaf,
		DeviceCount: info.DeviceCount,
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
