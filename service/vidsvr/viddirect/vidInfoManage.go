package viddirect

import (
	client "github.com/i-Things/things/service/vidsvr/client/vidmgrinfomanage"
	server "github.com/i-Things/things/service/vidsvr/internal/server/vidmgrinfomanage"
)

var (
	vidmgrInfoSvr client.VidmgrInfoManage
)

// 服务信息管理
func NewVidmgrManage(runSvr bool) client.VidmgrInfoManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectVidmgrInfoManage(svcCtx, server.NewVidmgrInfoManageServer(svcCtx))
	return vidSvr
}
