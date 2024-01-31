package dmdirect

import (
	client "github.com/i-Things/things/service/dmsvr/client/otamodulemanage"
	server "github.com/i-Things/things/service/dmsvr/internal/server/otamodulemanage"
)

var (
	otaModuleManageSvr client.OTAModuleManage
)

func NewOTAModuleManage(runSvr bool) client.OTAModuleManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectOTAModuleManage(svcCtx, server.NewOTAModuleManageServer(svcCtx))
	return dmSvr
}
