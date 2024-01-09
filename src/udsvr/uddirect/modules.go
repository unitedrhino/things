package uddirect

import (
	"github.com/i-Things/things/src/udsvr/client/rule"
	server "github.com/i-Things/things/src/udsvr/internal/server/rule"
)

func NewRule(runSvr bool) rule.Rule {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return rule.NewDirectRule(svcCtx, server.NewRuleServer(svcCtx))
}
