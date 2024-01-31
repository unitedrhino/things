package viddirect

import (
	client "github.com/i-Things/things/service/vidsvr/client/vidmgrconfigmanage"
	server "github.com/i-Things/things/service/vidsvr/internal/server/vidmgrconfigmanage"
)

var (
	vidmgrConfigSvr client.VidmgrConfigManage
)

// 服务配置管理
func NewVidmgrConfigManage(runSvr bool) client.VidmgrConfigManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrConfigManage(svcCtx, server.NewVidmgrConfigManageServer(svcCtx))
	return vidSvr
}
