package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/otafirmwaremanage"
	server "github.com/i-Things/things/src/dmsvr/internal/server/otafirmwaremanage"
)

var (
	otaFirmwareManageSvr client.OTAFirmwareManage
)

func NewOTAFirmwareManage(runSvr bool) client.OTAFirmwareManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectOTAFirmwareManage(svcCtx, server.NewOTAFirmwareManageServer(svcCtx))
	return dmSvr
}
