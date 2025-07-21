package pubDev

import (
	"context"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"time"
)

type (
	CompareMsg func(payload []byte) bool

	PubDev interface {
		PublishToDev(ctx context.Context, msg *deviceMsg.PublishMsg) error
		ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, timeout time.Duration, compareMsg CompareMsg) ([]byte, error)
	}
)

var s *protocol.ScriptTrans

func NewPubDev(fast *eventBus.FastEvent, S *protocol.ScriptTrans) (PubDev, error) {
	s = S
	return newPubDevClient(fast)
}
