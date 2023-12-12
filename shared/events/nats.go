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
	"strings"
	"time"
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
		msg.Ack()
		utils.Go(context.Background(), func() {
			var ctx context.Context
			utils.Recover(ctx)
			startTime := time.Now()
			emsg := GetEventMsg(msg.Data)
			if emsg == nil {
				logx.Error(msg.Subject, string(msg.Data))
				return
			}
			ctx = emsg.GetCtx()
			ctx, span := ctxs.StartSpan(ctx, msg.Subject, "")
			defer span.End()

			err := handle(ctx, emsg.GetData(), msg)
			duration := time.Now().Sub(startTime)
			if err != nil {
				logx.WithContext(ctx).WithDuration(duration).Errorf("nats subscription|startTime:%v,subject:%v,body:%v,err:%v",
					startTime, msg.Subject, string(emsg.GetData()), err)
			} else {
				logx.WithContext(ctx).WithDuration(duration).Infof("nats subscription|startTime:%v,subject:%v,body:%v",
					startTime, msg.Subject, string(emsg.GetData()))
			}
		})

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
