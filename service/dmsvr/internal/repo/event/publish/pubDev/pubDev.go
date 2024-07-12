package pubDev

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/eventBus"
	"time"
)

type (
	CompareMsg func(payload []byte) bool

	PubDev interface {
		PublishToDev(ctx context.Context, msg *deviceMsg.PublishMsg) error
		ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, timeout time.Duration, compareMsg CompareMsg) ([]byte, error)
	}
)

func NewPubDev(fast *eventBus.FastEvent) (PubDev, error) {
	return newNatsClient(fast)
}
