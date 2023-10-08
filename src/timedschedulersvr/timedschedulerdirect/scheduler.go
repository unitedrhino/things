package timedschedulerdirect

import (
	client "github.com/i-Things/things/src/timedschedulersvr/client/scheduler"
	server "github.com/i-Things/things/src/timedschedulersvr/internal/server/scheduler"
)

var (
	schedulerSvr client.Scheduler
)

func NewSchedulerMsg(runSvr bool) client.Scheduler {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectScheduler(svcCtx, server.NewSchedulerServer(svcCtx))
	return dmSvr
}
