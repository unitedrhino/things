package dataUpdate

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/events"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	nc, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: nc}, nil
}

func (n *NatsJsClient) UpdateWithTopic(ctx context.Context, topic string, info any) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(topic, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s info:%v,err:%v", utils.FuncName(),
		utils.Fmt(info), err)
	return err
}
