package dmdirect

import (
	client "github.com/i-Things/things/service/dmsvr/client/otaupgradetaskmanage"
	server "github.com/i-Things/things/service/dmsvr/internal/server/otaupgradetaskmanage"
)

var (
	otaUpgradeTaskManageSvr client.OTAUpgradeTaskManage
)

func NewOTAUpgradeTaskManage(runSvr bool) client.OTAUpgradeTaskManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectOTAUpgradeTaskManage(svcCtx, server.NewOTAUpgradeTaskManageServer(svcCtx))
	return dmSvr
}
