package alarm

import (
	"github.com/i-Things/things/shared/def"
)

type DealRecordFilter struct {
	AlarmRecordID int64 //告警配置ID
	Time          def.TimeRange
}

const (
	DealTypeHuman  = 1 //人工处理
	DealTypeSystem = 2 //系统处理
)
