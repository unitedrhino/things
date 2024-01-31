package alarm

import (
	"gitee.com/i-Things/core/shared/def"
)

type LogFilter struct {
	AlarmRecordID int64 //告警配置ID
	Time          def.TimeRange
}
