package ruledirect

import (
	client "github.com/i-Things/things/src/rulesvr/client/scenelinkage"
	server "github.com/i-Things/things/src/rulesvr/internal/server/scenelinkage"
)

func NewSceneLinkage() client.SceneLinkage {
	svc := GetCtxSvc()
	svr := client.NewDirectSceneLinkage(svc, server.NewSceneLinkageServer(svc))
	return svr
}
