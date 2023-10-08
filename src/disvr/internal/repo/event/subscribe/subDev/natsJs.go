package subDev

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceStatus"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

var (
	natsJsConsumerName = "disvr"
	natsJsConsumerPos  int
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	nc, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: nc}, nil
}

func (n *NatsJsClient) Subscribe(handle Handle) error {
	err := n.queueSubscribeDevPublish(topics.DeviceUpThingAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Thing(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.Thing failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpOtaAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Ota(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.Ota failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpConfigAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Config(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.Config failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpSDKLogAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).SDKLog(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.SDKLog failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpShadowAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Shadow(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.Shadow failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpGatewayAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Gateway(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.Gateway failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}

	_, err = n.client.QueueSubscribe(topics.DeviceUpStatusConnected, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ele, err := deviceStatus.GetDevConnMsg(ctx, msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.GetDevConnMsg failure err:%v", utils.FuncName(), err)
				return err
			}
			return handle(ctx).Connected(ele)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DeviceUpStatusConnected)))

	if err != nil {
		return err
	}

	_, err = n.client.QueueSubscribe(topics.DeviceUpStatusDisconnected, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ele, err := deviceStatus.GetDevConnMsg(ctx, msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.GetDevConnMsg failure err:%v", utils.FuncName(), err)
				return err
			}
			return handle(ctx).Disconnected(ele)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DeviceUpStatusDisconnected)))
	if err != nil {
		return err
	}

	err = n.queueSubscribeDevPublish(topics.DeviceUpExtAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Ext(msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.Ext failure err:%v", utils.FuncName(), err)
				return err
			}
			return err
		})
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsJsClient) queueSubscribeDevPublish(topic string,
	handleFunc func(ctx context.Context, msg *deviceMsg.PublishMsg) error) error {
	_, err := n.client.QueueSubscribe(topic, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			defer utils.Recover(ctx)
			ele, err := deviceMsg.GetDevPublish(ctx, msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.GetDevPublish failure err:%v", utils.FuncName(), err)
				return err
			}
			return handleFunc(ctx, ele)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topic)))
	if err != nil {
		return err
	}
	return nil
}
