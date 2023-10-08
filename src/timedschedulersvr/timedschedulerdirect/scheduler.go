package timedschedulerdirect

import (
	client "github.com/i-Things/things/src/timedschedulersvr/client/timedscheduler"
	server "github.com/i-Things/things/src/timedschedulersvr/internal/server/timedscheduler"
)

var (
	schedulerSvr client.Timedscheduler
)

func NewScheduler(runSvr bool) client.Timedscheduler {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := client.NewDirectTimedscheduler(svcCtx, server.NewTimedschedulerServer(svcCtx))
	return svr
}
