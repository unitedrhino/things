package timedjobdirect

import (
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
	server "github.com/i-Things/things/src/timed/timedjobsvr/internal/server/timedmanage"
)

var (
	jobSvr timedmanage.TimedManage
)

func NewTimedJob(runSvr bool) timedmanage.TimedManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := timedmanage.NewDirectTimedManage(svcCtx, server.NewTimedManageServer(svcCtx))
	return svr
}
