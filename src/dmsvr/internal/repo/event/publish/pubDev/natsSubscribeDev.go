package pubDev

import (
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	natsSubDev struct {
		subscription *nats.Subscription
	}
)

func newNatsSubDev(subscription *nats.Subscription) *natsSubDev {
	return &natsSubDev{subscription: subscription}
}

func (s *natsSubDev) UnSubscribe() error {
	return s.subscription.Unsubscribe()
}

func (s *natsSubDev) GetMsg(timeout time.Duration) (ele *device.PublishMsg, err error) {
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
	//向jaeger推送当前节点信息，路径名为主题名
	ctx, span := traces.StartSpan(ctx, msg.Subject, "")
	logx.Infof("[dm.GetMsg]|-------------------trace:%s, spanid:%s|topic:%s",
		span.SpanContext().TraceID(), span.SpanContext().SpanID(), msg.Subject)
	defer span.End()
	ele, err = device.GetDevPublish(ctx, emsg.GetData())
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return
	}
	return
}
