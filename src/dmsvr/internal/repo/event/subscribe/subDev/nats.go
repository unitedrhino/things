package subDev

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceStatus"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client *clients.NatsClient
	}
)

const (
	ThingsDeliverGroup = "things_dm_group"
	natsJsConsumerName = "dmsvr"
)

func newNatsClient(conf conf.EventConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, natsJsConsumerName, conf.Nats)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
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

	err = n.client.QueueSubscribe(topics.DeviceUpStatusConnected, ThingsDeliverGroup,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ele, err := deviceStatus.GetDevConnMsg(ctx, msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.GetDevConnMsg failure err:%v", utils.FuncName(), err)
				return err
			}
			return handle(ctx).Connected(ele)
		})

	if err != nil {
		return err
	}

	err = n.client.QueueSubscribe(topics.DeviceUpStatusDisconnected, ThingsDeliverGroup,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ele, err := deviceStatus.GetDevConnMsg(ctx, msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.GetDevConnMsg failure err:%v", utils.FuncName(), err)
				return err
			}
			return handle(ctx).Disconnected(ele)
		})
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

func (n *NatsClient) queueSubscribeDevPublish(topic string,
	handleFunc func(ctx context.Context, msg *deviceMsg.PublishMsg) error) error {
	err := n.client.QueueSubscribe(topic, ThingsDeliverGroup,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			defer utils.Recover(ctx)
			ele, err := deviceMsg.GetDevPublish(ctx, msg)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.GetDevPublish failure err:%v", utils.FuncName(), err)
				return err
			}
			return handleFunc(ctx, ele)
		})
	if err != nil {
		return err
	}
	return nil
}
