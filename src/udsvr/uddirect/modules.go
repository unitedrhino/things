package uddirect

import (
	"github.com/i-Things/things/src/udsvr/client/ops"
	"github.com/i-Things/things/src/udsvr/client/rule"
	opsServer "github.com/i-Things/things/src/udsvr/internal/server/ops"
	ruleServer "github.com/i-Things/things/src/udsvr/internal/server/rule"
)

func NewRule(runSvr bool) rule.Rule {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return rule.NewDirectRule(svcCtx, ruleServer.NewRuleServer(svcCtx))
}

func NewOps(runSvr bool) ops.Ops {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return ops.NewDirectOps(svcCtx, opsServer.NewOpsServer(svcCtx))
}
