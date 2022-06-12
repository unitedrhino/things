package events

import (
	"context"
	"github.com/i-Things/things/shared/traces"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type HandleFunc func(ctx context.Context, msg []byte) error

func NatsSubscription(handle HandleFunc) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		msg.Ack()
		emsg := GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ctx, span := traces.StartSpan(ctx, msg.Subject, "")
		logx.Infof("[dmsvr.NatsSubscription]|-------------------trace:%s, spanid:%s|topic:%s",
			span.SpanContext().TraceID(), span.SpanContext().SpanID(), msg.Subject)
		defer span.End()
		err := handle(ctx, emsg.GetData())
		logx.WithContext(ctx).Infof("nats subscription|subject:%v,data:%v,err:%v",
			msg.Subject, string(msg.Data), err)
	}
}
