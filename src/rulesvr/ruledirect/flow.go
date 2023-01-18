package ruledirect

import (
	client "github.com/i-Things/things/src/rulesvr/client/flow"
	server "github.com/i-Things/things/src/rulesvr/internal/server/flow"
)

var (
	deviceManageSvr client.Flow
)

func NewFlow() client.Flow {
	svc := GetCtxSvc()
	svr := client.NewDirectFlow(svc, server.NewFlowServer(svc))
	return svr
}
