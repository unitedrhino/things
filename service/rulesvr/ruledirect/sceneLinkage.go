package ruledirect

import (
	client "github.com/i-Things/things/service/rulesvr/client/scenelinkage"
	server "github.com/i-Things/things/service/rulesvr/internal/server/scenelinkage"
)

func NewSceneLinkage(runSvr bool) client.SceneLinkage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := client.NewDirectSceneLinkage(svcCtx, server.NewSceneLinkageServer(svcCtx))
	return svr
}
