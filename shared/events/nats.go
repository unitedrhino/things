package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/core/timex"
	"strings"
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
		var ctx context.Context
		utils.Recover(ctx)
		startTime := timex.Now()
		msg.Ack()
		emsg := GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx = emsg.GetCtx()
		ctx, span := ctxs.StartSpan(ctx, msg.Subject, "")
		defer span.End()
		err := handle(ctx, emsg.GetData(), msg)
		duration := timex.Since(startTime)
		if err != nil {
			logx.WithContext(ctx).WithDuration(duration).Errorf("nats subscription|subject:%v,body:%v,err:%v",
				msg.Subject, string(emsg.GetData()), err)
		} else {
			logx.WithContext(ctx).WithDuration(duration).Infof("nats subscription|subject:%v,body:%v",
				msg.Subject, string(emsg.GetData()))
		}
	}
}
func GenNatsJsDurable(serverName string, topic string) string {
	ip := netx.InternalIp()
	ret := fmt.Sprintf("%s_%s_%s", serverName, ip, topic)
	ret = strings.ReplaceAll(ret, ".", "-")
	ret = strings.ReplaceAll(ret, "*", "+")
	ret = strings.ReplaceAll(ret, ">", "~")
	return ret
}
