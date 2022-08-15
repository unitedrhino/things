package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/productmanage"
	server "github.com/i-Things/things/src/dmsvr/internal/server/productmanage"
)

var (
	productManageSvr client.ProductManage
)

func NewProductManage(config *Config) client.ProductManage {
	svc := getCtxSvc(config)
	dmSvr := client.NewDirectProductManage(svc, server.NewProductManageServer(svc))
	return dmSvr
}
