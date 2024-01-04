package viddirect

import (
	client "github.com/i-Things/things/src/vidsvr/client/vidmgrgbsipmanage"
	server "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrgbsipmanage"
)

var (
	vidmgrGbsipInfo client.VidmgrGbsipManage
)

func NewVidmgrGbsipManage(runSvr bool) client.VidmgrGbsipManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrGbsipManage(svcCtx, server.NewVidmgrGbsipManageServer(svcCtx))
	return vidSvr
}
