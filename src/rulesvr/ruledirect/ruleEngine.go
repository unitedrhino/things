package ruledirect

import (
	client "github.com/i-Things/things/src/rulesvr/client/ruleengine"
	server "github.com/i-Things/things/src/rulesvr/internal/server/ruleengine"
)

func NewRuleEngine() client.RuleEngine {
	svc := GetCtxSvc()
	svr := client.NewDirectRuleEngine(svc, server.NewRuleEngineServer(svc))
	return svr
}
