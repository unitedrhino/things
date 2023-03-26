package alarmcenterlogic

import (
	"github.com/i-Things/things/shared/def"
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
		CreatedTime: in.CreatedTime.Unix(),
	}
}
func ToTimeRange(timeRange *rule.TimeRange) def.TimeRange {
	if timeRange == nil {
		return def.TimeRange{}
	}
	return def.TimeRange{Start: timeRange.Start, End: timeRange.End}
}
func ToAlarmDealRecord(in *mysql.RuleAlarmDealRecord) *rule.AlarmDeal {
	return &rule.AlarmDeal{
		Id:            in.Id,
		AlarmRecordID: in.AlarmRecordID,
		Result:        in.Result,
		Type:          in.Type,
		AlarmTime:     utils.TimeToInt64(in.AlarmTime),
		CreatedTime:   utils.TimeToInt64(in.CreatedTime),
	}
}
func ToAlarmLog(in *mysql.RuleAlarmLog) *rule.AlarmLog {
	return &rule.AlarmLog{
		Id:            in.Id,
		AlarmRecordID: in.AlarmRecordID,
		Serial:        in.Serial,
		SceneName:     in.SceneName,
		SceneID:       in.SceneID,
		Desc:          in.Desc,
		CreatedTime:   utils.TimeToInt64(in.CreatedTime),
	}
}
func ToAlarmRecord(in *mysql.RuleAlarmRecord) *rule.AlarmRecord {
	return &rule.AlarmRecord{
		Id:          in.Id,
		AlarmID:     in.AlarmID,
		TriggerType: in.TriggerType,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		SceneName:   in.SceneName,
		SceneID:     in.SceneID,
		Level:       in.Level,
		LastAlarm:   utils.TimeToInt64(in.LastAlarm),
		CreatedTime: utils.TimeToInt64(in.CreatedTime),
	}
}
