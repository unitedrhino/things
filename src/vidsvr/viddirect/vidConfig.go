package viddirect

import (
	client "github.com/i-Things/things/src/vidsvr/client/vidmgrconfigmange"
	server "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrconfigmange"
)

var (
	vidmgrConfigSvr client.VidmgrConfigMange
)

func NewVidmgrConfigManage(runSvr bool) client.VidmgrConfigMange {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrConfigMange(svcCtx, server.NewVidmgrConfigMangeServer(svcCtx))
	return vidSvr
}
