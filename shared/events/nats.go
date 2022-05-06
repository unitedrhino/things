package events

import (
	"context"
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
		err := handle(ctx, emsg.GetData())
		logx.WithContext(ctx).Infof("nats subscription|subject:%v,data:%v,err:%v",
			msg.Subject, string(msg.Data), err)
	}
}
