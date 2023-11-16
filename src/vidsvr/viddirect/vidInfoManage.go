package viddirect

import (
	client "github.com/i-Things/things/src/vidsvr/client/vidmgrinfomanage"
	server "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrinfomanage"
)

var (
	vidmgrInfoSvr client.VidmgrInfoManage
)

func NewVidmgrManage(runSvr bool) client.VidmgrInfoManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrInfoManage(svcCtx, server.NewVidmgrInfoManageServer(svcCtx))
	//dmSvr := client.NewDirectProductManage(svcCtx, server.NewProductManageServer(svcCtx))
	return vidSvr
}
