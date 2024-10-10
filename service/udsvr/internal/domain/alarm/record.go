package alarm

import "gitee.com/unitedrhino/share/def"

type RecordFilter struct {
	AlarmID     int64 // 告警配置ID
	TriggerType int64
	ProductID   string
	DeviceName  string
	Time        def.TimeRange
}
