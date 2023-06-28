package events

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/traces"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

type HandleFunc func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error

func NatsSubWithType[msgType any](handle func(ctx context.Context, msgIn msgType, natsMsg *nats.Msg) error) func(msg *nats.Msg) {
	return NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
		var tempInfo msgType
		err := json.Unmarshal(msg, &tempInfo)
		if err != nil {
			return err
		}
		return handle(ctx, tempInfo, natsMsg)
	})
}

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
		err := handle(ctx, emsg.GetData(), msg)
		if err != nil {
			logx.WithContext(ctx).WithDuration(duration).Errorf("nats subscription|subject:%v,body:%v,err:%v",
				msg.Subject, string(emsg.GetData()), err)
		} else {
			logx.WithContext(ctx).WithDuration(duration).Infof("nats subscription|subject:%v",
				msg.Subject)
		}
	}
}
