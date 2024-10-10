package dmdirect

import (
	client "gitee.com/unitedrhino/things/service/dmsvr/client/otamanage"
	server "gitee.com/unitedrhino/things/service/dmsvr/internal/server/otamanage"
)

var (
	otaFirmwareManageSvr client.OtaManage
)

func NewOtaManage(runSvr bool) client.OtaManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectOtaManage(svcCtx, server.NewOtaManageServer(svcCtx))
	return dmSvr
}
