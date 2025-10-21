package pubDev

import (
	"context"
	"time"

	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
)

type (
	CompareMsg func(payload []byte) bool

	PubDev interface {
		PublishToDev(ctx context.Context, di *dm.DeviceInfo, msg *deviceMsg.PublishMsg) error
		ReqToDeviceSync(ctx context.Context, di *dm.DeviceInfo, reqMsg *deviceMsg.PublishMsg, timeout time.Duration, compareMsg CompareMsg) ([]byte, error)
	}
)

var s *protocol.ScriptTrans

func NewPubDev(fast *eventBus.FastEvent, S *protocol.ScriptTrans) (PubDev, error) {
	s = S
	return newPubDevClient(fast)
}
