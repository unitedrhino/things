package alarm

import (
	"github.com/i-Things/things/shared/def"
)

type LogFilter struct {
	AlarmRecordID int64 //告警配置ID
	Time          def.TimeRange
}
