package msgExt

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"time"
)

type (
	Resp struct {
		*deviceMsg.CommonMsg
		DeviceSendTime int64 `json:"deviceSendTime"` //ntp设备发送毫秒时间戳
		ServerRecvTime int64 `json:"serverRecvTime"` //ntp云端接收毫秒时间戳
	}
)

func (d *Resp) GetTimeStamp(defaultTime time.Time) time.Time {
	if d.Timestamp == 0 {
		return defaultTime
	}
	return time.UnixMilli(d.Timestamp)
}
