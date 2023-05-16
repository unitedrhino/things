package ruledirect

import (
	client "github.com/i-Things/things/src/rulesvr/client/ruleengine"
	server "github.com/i-Things/things/src/rulesvr/internal/server/ruleengine"
)

func NewRuleEngine(runSvr bool) client.RuleEngine {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := client.NewDirectRuleEngine(svcCtx, server.NewRuleEngineServer(svcCtx))
	return svr
}
