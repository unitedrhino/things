package alarmcenterlogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
	"time"
)

func ToAlarmInfoPo(in *rule.AlarmInfo) *mysql.RuleAlarmInfo {
	return &mysql.RuleAlarmInfo{
		Id:          in.Id,
		Name:        in.Name,
		Desc:        in.Desc,
		Level:       in.Level,
		State:       in.State,
		DealState:   in.DealState,
		LastAlarm:   utils.ToNullTime(in.LastAlarm),
		CreatedTime: time.Unix(in.CreatedTime, 0),
	}
}
func ToAlarmInfo(in *mysql.RuleAlarmInfo) *rule.AlarmInfo {
	return &rule.AlarmInfo{
		Id:          in.Id,
		Name:        in.Name,
		Desc:        in.Desc,
		Level:       in.Level,
		State:       in.State,
		DealState:   in.DealState,
		LastAlarm:   utils.GetNullTime(in.LastAlarm),
		CreatedTime: in.CreatedTime.Unix(),
	}
}
