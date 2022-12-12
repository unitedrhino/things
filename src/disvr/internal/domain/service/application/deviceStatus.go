package application

import (
	"github.com/i-Things/things/shared/devices"
	"time"
)

//连接和断连消息信息
type ConnectMsg struct {
	Device    devices.Core
	Timestamp time.Time
}
