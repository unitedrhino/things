package uddirect

import (
	"github.com/i-Things/things/service/udsvr/client/ops"
	"github.com/i-Things/things/service/udsvr/client/rule"
	"github.com/i-Things/things/service/udsvr/client/userdevice"
	opsServer "github.com/i-Things/things/service/udsvr/internal/server/ops"
	ruleServer "github.com/i-Things/things/service/udsvr/internal/server/rule"
	userdeviceServer "github.com/i-Things/things/service/udsvr/internal/server/userdevice"
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
func NewUserDevice(runSvr bool) userdevice.UserDevice {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return userdevice.NewDirectUserDevice(svcCtx, userdeviceServer.NewUserDeviceServer(svcCtx))
}
