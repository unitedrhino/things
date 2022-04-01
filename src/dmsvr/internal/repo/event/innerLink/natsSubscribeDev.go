package innerLink

import (
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	SubDev struct {
		subscription *nats.Subscription
	}
)

func NewSubDev(subscription *nats.Subscription) *SubDev {
	return &SubDev{subscription: subscription}
}

func (s *SubDev) UnSubscribe() error {
	return s.subscription.Unsubscribe()
}

func (s *SubDev) GetMsg(timeout time.Duration) (ele *deviceMsg.Elements, err error) {
	msg, err := s.subscription.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	msg.Ack()
	emsg := events.GetEventMsg(msg.Data)
	if emsg == nil {
		logx.Error(msg.Subject, string(msg.Data))
		return
	}
	ctx := emsg.GetCtx()
	ele, err = deviceMsg.GetDevPublish(ctx, emsg.GetData())
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return
	}
	return
}
