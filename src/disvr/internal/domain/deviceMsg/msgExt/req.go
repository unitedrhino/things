package msgExt

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"time"
)

type (
	Req struct {
		deviceMsg.CommonMsg
	}
)

func (d *Req) GetTimeStamp(defaultTime int64) time.Time {
	if d.Timestamp == 0 {
		return time.UnixMilli(defaultTime)
	}
	return time.UnixMilli(d.Timestamp)
}
