package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/otajobmanage"
	server "github.com/i-Things/things/src/dmsvr/internal/server/otajobmanage"
)

var (
	otaJobManageSvr client.OTAJobManage
)

func NewOTAJobManage(runSvr bool) client.OTAJobManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectOTAJobManage(svcCtx, server.NewOTAJobManageServer(svcCtx))
	return dmSvr
}
