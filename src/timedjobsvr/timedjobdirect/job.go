package timedjobdirect

import (
	client "github.com/i-Things/things/src/timedjobsvr/client/job"
	server "github.com/i-Things/things/src/timedjobsvr/internal/server/job"
)

var (
	jobSvr client.Job
)

func NewJob(runSvr bool) client.Job {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectJob(svcCtx, server.NewJobServer(svcCtx))
	return dmSvr
}
