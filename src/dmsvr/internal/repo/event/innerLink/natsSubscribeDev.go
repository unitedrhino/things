package innerLink

import (
	"github.com/i-Things/things/shared/events"
	my_trace "github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
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

func (s *SubDev) GetMsg(timeout time.Duration) (ele *device.PublishMsg, err error) {
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
	ctx, span := my_trace.StartSpan(ctx, msg.Subject, "")
	logx.Infof("[mqtt.GetMsg]|-------------------trace:%s, spanid:%s|topic:%s",
		span.SpanContext().TraceID(), span.SpanContext().SpanID(), msg.Subject)
	defer span.End()
	ele, err = device.GetDevPublish(ctx, emsg.GetData())
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return
	}

	return
}
