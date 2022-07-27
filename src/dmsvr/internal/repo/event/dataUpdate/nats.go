package dataUpdate

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

const (
	TopicUpdateTempModel = "server.dm.update.tempModel"
)

func NewNatsClient(conf conf.NatsConf) (*NatsClient, error) {
	connectOpts := nats.Options{
		Url:      conf.Url,
		User:     conf.User,
		Password: conf.Pass,
		Token:    conf.Token,
	}
	nc, err := connectOpts.Connect()
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) TempModelUpdate(ctx context.Context, info *schema.SchemaInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = n.client.Publish(TopicUpdateTempModel, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s|info:%v,err:%v", utils.FuncName(),
		info, err)
	return err
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.Subscribe(TopicUpdateTempModel,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			tempInfo := schema.SchemaInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).TempModelClearCache(&tempInfo)
		}))
	if err != nil {
		return err
	}
	return nil
}
