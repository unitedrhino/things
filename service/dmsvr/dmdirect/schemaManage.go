package dmdirect

import (
	client "github.com/i-Things/things/service/dmsvr/client/schemamanage"
	server "github.com/i-Things/things/service/dmsvr/internal/server/schemamanage"
)

var (
	schemaManageSvr client.SchemaManage
)

func NewSchemaManage(runSvr bool) client.SchemaManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectSchemaManage(svcCtx, server.NewSchemaManageServer(svcCtx))
	return dmSvr
}
