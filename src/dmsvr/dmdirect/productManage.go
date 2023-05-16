package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/productmanage"
	server "github.com/i-Things/things/src/dmsvr/internal/server/productmanage"
)

var (
	productManageSvr client.ProductManage
)

func NewProductManage(runSvr bool) client.ProductManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectProductManage(svcCtx, server.NewProductManageServer(svcCtx))
	return dmSvr
}
