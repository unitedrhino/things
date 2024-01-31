package dmdirect

import (
	client "github.com/i-Things/things/service/dmsvr/client/firmwaremanage"
	server "github.com/i-Things/things/service/dmsvr/internal/server/firmwaremanage"
)

var (
	firmwareManageSvr client.FirmwareManage
)

func NewFirmwareManage(runSvr bool) client.FirmwareManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectFirmwareManage(svcCtx, server.NewFirmwareManageServer(svcCtx))
	return dmSvr
}
