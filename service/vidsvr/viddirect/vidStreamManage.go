package viddirect

import (
	client "github.com/i-Things/things/service/vidsvr/client/vidmgrstreammanage"
	server "github.com/i-Things/things/service/vidsvr/internal/server/vidmgrstreammanage"
)

// 视频流管理
var (
	vidmgrStreamSvr client.VidmgrStreamManage
)

func NewVidmgrStreamManage(runSvr bool) client.VidmgrStreamManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrStreamManage(svcCtx, server.NewVidmgrStreamManageServer(svcCtx))
	return vidSvr
}
