package dmdirect

import (
	client "gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	server "gitee.com/unitedrhino/things/service/dmsvr/internal/server/productmanage"
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
