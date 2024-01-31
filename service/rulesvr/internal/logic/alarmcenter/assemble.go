package alarmcenterlogic

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/rulesvr/pb/rule"
)

func ToAlarmInfoPo(in *rule.AlarmInfo) *relationDB.RuleAlarmInfo {
	return &relationDB.RuleAlarmInfo{
		Name:   in.Name,
		Desc:   in.Desc,
		Level:  in.Level,
		Status: in.Status,
	}
}
func ToAlarmInfo(in *relationDB.RuleAlarmInfo) *rule.AlarmInfo {
	return &rule.AlarmInfo{
		Id:          in.ID,
		Name:        in.Name,
		Desc:        in.Desc,
		Level:       in.Level,
		Status:      in.Status,
		CreatedTime: in.CreatedTime.Unix(),
	}
}
func ToTimeRange(timeRange *rule.TimeRange) def.TimeRange {
	if timeRange == nil {
		return def.TimeRange{}
	}
	return def.TimeRange{Start: timeRange.Start, End: timeRange.End}
}
func ToAlarmDealRecord(in *relationDB.RuleAlarmDealRecord) *rule.AlarmDeal {
	return &rule.AlarmDeal{
		Id:            in.ID,
		AlarmRecordID: in.AlarmRecordID,
		Result:        in.Result,
		Type:          in.Type,
		AlarmTime:     utils.TimeToInt64(in.AlarmTime),
		CreatedTime:   utils.TimeToInt64(in.CreatedTime),
	}
}
func ToAlarmLog(in *relationDB.RuleAlarmLog) *rule.AlarmLog {
	return &rule.AlarmLog{
		Id:            in.ID,
		AlarmRecordID: in.AlarmRecordID,
		Serial:        in.Serial,
		SceneName:     in.SceneName,
		SceneID:       in.SceneID,
		Desc:          in.Desc,
		CreatedTime:   utils.TimeToInt64(in.CreatedTime),
	}
}
func ToAlarmRecord(in *relationDB.RuleAlarmRecord) *rule.AlarmRecord {
	return &rule.AlarmRecord{
		Id:          in.ID,
		AlarmID:     in.AlarmID,
		TriggerType: in.TriggerType,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		SceneName:   in.SceneName,
		SceneID:     in.SceneID,
		Level:       in.Level,
		DealState:   in.DealState,
		LastAlarm:   utils.TimeToInt64(in.LastAlarm),
		CreatedTime: utils.TimeToInt64(in.CreatedTime),
	}
}
