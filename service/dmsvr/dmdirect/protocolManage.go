package dmdirect

import (
	client "github.com/i-Things/things/service/dmsvr/client/protocolmanage"
	server "github.com/i-Things/things/service/dmsvr/internal/server/protocolmanage"
)

var (
	protocolManageSvr client.ProtocolManage
)

func NewProtocolManage(runSvr bool) client.ProtocolManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectProtocolManage(svcCtx, server.NewProtocolManageServer(svcCtx))
	return dmSvr
}
