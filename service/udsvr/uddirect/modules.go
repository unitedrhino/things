package uddirect

import (
	"gitee.com/unitedrhino/things/service/udsvr/client/rule"
	ruleServer "gitee.com/unitedrhino/things/service/udsvr/internal/server/rule"
)

func NewRule(runSvr bool) rule.Rule {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return rule.NewDirectRule(svcCtx, ruleServer.NewRuleServer(svcCtx))
}
