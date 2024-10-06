package uddirect

import (
	"gitee.com/i-Things/things/service/udsvr/client/rule"
	ruleServer "gitee.com/i-Things/things/service/udsvr/internal/server/rule"
)

func NewRule(runSvr bool) rule.Rule {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return rule.NewDirectRule(svcCtx, ruleServer.NewRuleServer(svcCtx))
}
