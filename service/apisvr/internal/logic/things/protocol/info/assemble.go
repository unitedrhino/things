package info

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToInfoPb(in *types.ProtocolInfo) *dm.ProtocolInfo {
	if in == nil {
		return nil
	}
	return &dm.ProtocolInfo{
		Id:            in.ID,
		Name:          in.Name,
		Code:          in.Code,
		TransProtocol: in.TransProtocol,
		Desc:          in.Desc,
		Endpoints:     in.Endpoints,
		EtcdKey:       in.EtcdKey,
		ConfigFields:  ToConfigFieldsPb(in.ConfigFields),
		ConfigInfos:   ToConfigInfosPb(in.ConfigInfos),
	}
}

func ToConfigFieldsPb(in []*types.ProtocolConfigField) (ret []*dm.ProtocolConfigField) {
	for _, v := range in {
		ret = append(ret, &dm.ProtocolConfigField{
			Id:         v.ID,
			Group:      v.Group,
			Key:        v.Key,
			Label:      v.Label,
			IsRequired: v.IsRequired,
			Sort:       v.Sort,
		})
	}
	return
}

func ToConfigInfosPb(in []*types.ProtocolConfigInfo) (ret []*dm.ProtocolConfigInfo) {
	for _, v := range in {
		ret = append(ret, &dm.ProtocolConfigInfo{
			Id:     v.ID,
			Config: v.Config,
			Desc:   v.Desc,
		})
	}
	return
}

func ToInfoTypes(in *dm.ProtocolInfo) *types.ProtocolInfo {
	if in == nil {
		return nil
	}
	return &types.ProtocolInfo{
		ID:            in.Id,
		Name:          in.Name,
		Code:          in.Code,
		TransProtocol: in.TransProtocol,
		Desc:          in.Desc,
		Endpoints:     in.Endpoints,
		EtcdKey:       in.EtcdKey,
		ConfigFields:  ToConfigFieldsTypes(in.ConfigFields),
		ConfigInfos:   ToConfigInfosTypes(in.ConfigInfos),
	}
}
func ToInfosTypes(in []*dm.ProtocolInfo) (ret []*types.ProtocolInfo) {
	for _, v := range in {
		ret = append(ret, ToInfoTypes(v))
	}
	return
}

func ToConfigFieldsTypes(in []*dm.ProtocolConfigField) (ret []*types.ProtocolConfigField) {
	for _, v := range in {
		ret = append(ret, &types.ProtocolConfigField{
			ID:         v.Id,
			Group:      v.Group,
			Key:        v.Key,
			Label:      v.Label,
			IsRequired: v.IsRequired,
			Sort:       v.Sort,
		})
	}
	return
}

func ToConfigInfosTypes(in []*dm.ProtocolConfigInfo) (ret []*types.ProtocolConfigInfo) {
	for _, v := range in {
		ret = append(ret, &types.ProtocolConfigInfo{
			ID:     v.Id,
			Config: v.Config,
			Desc:   v.Desc,
		})
	}
	return
}
