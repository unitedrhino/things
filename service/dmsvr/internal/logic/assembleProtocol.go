package logic

import (
	"github.com/i-Things/things/service/dmsvr/internal/domain/protocol"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToProtocolInfoPb(in *relationDB.DmProtocolInfo) *dm.ProtocolInfo {
	if in == nil {
		return nil
	}
	dpi := &dm.ProtocolInfo{
		Id:            in.ID,
		Name:          in.Name,
		Code:          in.Code,
		TransProtocol: in.TransProtocol,
		ConfigFields:  ToProtocolConfigFieldsPb(in.ConfigFields),
		ConfigInfos:   ToProtocolConfigInfosPb(in.ConfigInfos),
		Desc:          in.Desc,
		Endpoints:     in.Endpoints,
		EtcdKey:       in.EtcdKey,
	}

	return dpi
}

func ToProtocolConfigInfosPb(in protocol.ConfigInfos) (ret []*dm.ProtocolConfigInfo) {
	for _, v := range in {
		ret = append(ret, &dm.ProtocolConfigInfo{
			Id:     v.ID,
			Config: v.Config,
			Desc:   v.Desc,
		})
	}
	return
}

func ToProtocolConfigFieldsPb(in protocol.ConfigFields) (ret []*dm.ProtocolConfigField) {
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

func ToProtocolInfoPo(in *dm.ProtocolInfo) *relationDB.DmProtocolInfo {
	if in == nil {
		return nil
	}
	dpi := &relationDB.DmProtocolInfo{
		ID:            in.Id,
		Name:          in.Name,
		Code:          in.Code,
		TransProtocol: in.TransProtocol,
		ConfigFields:  ToProtocolConfigFieldsPo(in.ConfigFields),
		ConfigInfos:   ToProtocolConfigInfosPo(in.ConfigInfos),
		Desc:          in.Desc,
		Endpoints:     in.Endpoints,
		EtcdKey:       in.EtcdKey,
	}

	return dpi
}

func ToProtocolConfigInfosPo(in []*dm.ProtocolConfigInfo) (ret protocol.ConfigInfos) {
	for _, v := range in {
		ret = append(ret, &protocol.ConfigInfo{
			ID:     v.Id,
			Config: v.Config,
			Desc:   v.Desc,
		})
	}
	return
}

func ToProtocolConfigFieldsPo(in []*dm.ProtocolConfigField) (ret protocol.ConfigFields) {
	for _, v := range in {
		ret = append(ret, &protocol.ConfigField{
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
