package viddirect

import (
	client "github.com/i-Things/things/src/vidsvr/client/vidmgrmange"
	server "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrmange"
)

var (
	vidmgrMangeSvr client.VidmgrMange
)

func NewVidmgrManage(runSvr bool) client.VidmgrMange {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrMange(svcCtx, server.NewVidmgrMangeServer(svcCtx))
	//dmSvr := client.NewDirectProductManage(svcCtx, server.NewProductManageServer(svcCtx))
	return vidSvr
}
