package ruledirect

import (
	client "github.com/i-Things/things/src/rulesvr/client/alarmcenter"
	server "github.com/i-Things/things/src/rulesvr/internal/server/alarmcenter"
)

func NewAlarmCenter(runSvr bool) client.AlarmCenter {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	svr := client.NewDirectAlarmCenter(svcCtx, server.NewAlarmCenterServer(svcCtx))
	return svr
}
