package alarm

import "github.com/i-Things/things/shared/def"

type RecordFilter struct {
	AlarmID     int64 // 告警配置ID
	TriggerType int64
	ProductID   string
	DeviceName  string
	Time        def.TimeRange
}
