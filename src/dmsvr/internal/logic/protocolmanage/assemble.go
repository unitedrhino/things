package protocolmanagelogic

import (
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToProtocolInfo(in *relationDB.DmProtocolInfo) *dm.ProtocolInfo {
	dpi := &dm.ProtocolInfo{
		Id:           in.ID,
		Name:         in.Name,
		Protocol:     in.Protocol,
		ProtocolType: in.ProtocolType,
		Desc:         in.Desc,
		Endpoints:    in.Endpoints,
		EtcdKey:      in.EtcdKey,
	}

	return dpi
}
