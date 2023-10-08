package timedjobdirect

import (
	client "github.com/i-Things/things/src/timedjobsvr/client/timedjob"
	server "github.com/i-Things/things/src/timedjobsvr/internal/server/timedjob"
)

var (
	jobSvr client.TimedJob
)

func NewTimedJob(runSvr bool) client.TimedJob {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := client.NewDirectTimedJob(svcCtx, server.NewTimedJobServer(svcCtx))
	return svr
}
