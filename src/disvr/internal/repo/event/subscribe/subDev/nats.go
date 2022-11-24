package subDev

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceStatus"
	"github.com/nats-io/nats.go"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

const (
	ThingsDeliverGroup = "things_dm_group"
)

func newNatsClient(conf conf.NatsConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	err := n.queueSubscribeDevPublish(topics.DeviceUpThingAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Thing(msg)
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpOtaAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Ota(msg)
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpConfigAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Config(msg)
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpSDKLogAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).SDKLog(msg)
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpShadowAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Shadow(msg)
			return err
		})
	if err != nil {
		return err
	}
	err = n.queueSubscribeDevPublish(topics.DeviceUpGatewayAll,
		func(ctx context.Context, msg *deviceMsg.PublishMsg) error {
			err := handle(ctx).Gateway(msg)
			return err
		})
	if err != nil {
		return err
	}

	_, err = n.client.QueueSubscribe(topics.DeviceUpStatusConnected, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := deviceStatus.GetDevConnMsg(ctx, msg)
			if err != nil {
				return err
			}
			return handle(ctx).Connected(ele)
		}))

	if err != nil {
		return err
	}

	_, err = n.client.QueueSubscribe(topics.DeviceUpStatusDisconnected, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := deviceStatus.GetDevConnMsg(ctx, msg)
			if err != nil {
				return err
			}
			return handle(ctx).Disconnected(ele)
		}))
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsClient) queueSubscribeDevPublish(topic string,
	handleFunc func(ctx context.Context, msg *deviceMsg.PublishMsg) error) error {
	_, err := n.client.QueueSubscribe(topic, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := deviceMsg.GetDevPublish(ctx, msg)
			if err != nil {
				return err
			}
			return handleFunc(ctx, ele)
		}))
	if err != nil {
		return err
	}
	return nil
}
