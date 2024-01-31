package dmdirect

import (
	client "github.com/i-Things/things/service/dmsvr/client/otataskmanage"
	server "github.com/i-Things/things/service/dmsvr/internal/server/otataskmanage"
)

var (
	otaTaskManageSvr client.OtaTaskManage
)

func NewOtaTaskManage(runSvr bool) client.OtaTaskManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectOtaTaskManage(svcCtx, server.NewOtaTaskManageServer(svcCtx))
	return dmSvr
}
