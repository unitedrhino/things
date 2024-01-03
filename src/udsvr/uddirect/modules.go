package uddirect

import (
	"github.com/i-Things/things/src/udsvr/client/intelligentcontrol"
	server "github.com/i-Things/things/src/udsvr/internal/server/intelligentcontrol"
)

func NewIntelligentControl(runSvr bool) intelligentcontrol.IntelligentControl {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return intelligentcontrol.NewDirectIntelligentControl(svcCtx, server.NewIntelligentControlServer(svcCtx))
}
