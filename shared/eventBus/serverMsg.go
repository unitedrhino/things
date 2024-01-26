package eventBus

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
服务消息,不需要模糊匹配的,发送给所有订阅者的可以用这个来简化实现
*/

type ServerMsg struct {
	natsCli  *clients.NatsClient
	handlers map[string][]ServerFunc
}

type ServerFunc func(ctx context.Context, body []byte) error

func NewServerMsg(c conf.EventConf, serverName string) (s *ServerMsg, err error) {
	serverMsg := ServerMsg{handlers: map[string][]ServerFunc{}}
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		serverMsg.natsCli, err = clients.NewNatsClient2(c.Mode, serverName, c.Nats)
	default:
		err = errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
	}
	return &serverMsg, err
}
func (bus *ServerMsg) Start() error {
	for topic, handles := range bus.handlers {
		err := bus.natsCli.Subscribe(topic, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ctx = utils.CopyContext(ctx)
			for _, f := range handles {
				utils.Go(ctx, func() {
					err := f(ctx, msg)
					if err != nil {
						logx.WithContext(ctx).Error(err)
					}
				})
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅
func (bus *ServerMsg) Subscribe(topic string, f ServerFunc) {
	handler, ok := bus.handlers[topic]
	if !ok {
		handler = []ServerFunc{}
	}
	handler = append(handler, f)
	bus.handlers[topic] = handler
	return
}

// Publish 发布
// 这里异步执行，并且不会等待返回结果
func (bus *ServerMsg) Publish(ctx context.Context, topic string, arg any) error {
	err := bus.natsCli.Publish(topic, []byte(utils.Fmt(arg)))
	return err
}
