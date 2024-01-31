package alarm

import (
	"gitee.com/i-Things/share/def"
)

type LogFilter struct {
	AlarmRecordID int64 //告警配置ID
	Time          def.TimeRange
}
