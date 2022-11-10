package msgGateway

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"time"
)

type (
	GateWayReq struct {
		deviceMsg.CommonMsg
	}
)

func (d *GateWayReq) GetTimeStamp(logTime int64) time.Time {
	if logTime == 0 {
		if d.Timestamp != 0 {
			return time.UnixMilli(d.Timestamp)
		}
		return time.Now()
	}
	return time.UnixMilli(logTime)
}
