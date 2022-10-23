package events

import (
	"context"
	"github.com/i-Things/things/shared/traces"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

type HandleFunc func(ctx context.Context, msg []byte) error

func NatsSubscription(handle HandleFunc) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		startTime := timex.Now()
		duration := timex.Since(startTime)
		msg.Ack()
		emsg := GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ctx, span := traces.StartSpan(ctx, msg.Subject, "")
		defer span.End()
		err := handle(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).WithDuration(duration).Errorf("nats subscription|subject:%v,body:%v,err:%v",
				msg.Subject, string(emsg.GetData()), err)
		} else {
			logx.WithContext(ctx).WithDuration(duration).Infof("nats subscription|subject:%v",
				msg.Subject)
		}
	}
}
