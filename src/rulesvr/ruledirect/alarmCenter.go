package ruledirect

import (
	client "github.com/i-Things/things/src/rulesvr/client/alarmcenter"
	server "github.com/i-Things/things/src/rulesvr/internal/server/alarmcenter"
)

func NewAlarmCenter() client.AlarmCenter {
	svc := GetSvcCtx()
	svr := client.NewDirectAlarmCenter(svc, server.NewAlarmCenterServer(svc))
	return svr
}
