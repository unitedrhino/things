package viddirect

import (
	client "github.com/i-Things/things/src/vidsvr/client/vidmgrstreammanage"
	server "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrstreammanage"
)

var (
	vidmgrStreamSvr client.VidmgrStreamManage
)

func NewVidmgrStreamManage(runSvr bool) client.VidmgrStreamManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrStreamManage(svcCtx, server.NewVidmgrStreamManageServer(svcCtx))
	//dmSvr := client.NewDirectProductManage(svcCtx, server.NewProductManageServer(svcCtx))
	return vidSvr
}
