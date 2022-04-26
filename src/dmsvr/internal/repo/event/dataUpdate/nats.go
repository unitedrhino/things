package dataUpdate

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/templateModel"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client nats.JetStreamContext
	}
)

const (
	DmUpdateConsumeName  = "dm_rpc_update_consume"
	DmUpdateStreamName   = "dm_rpc_update_msg"
	TopicUpdate          = "dm.update"
	DmUpdateDeliverGroup = "dm_rpc_update_group"
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
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name: DmUpdateStreamName,
		Subjects: []string{
			TopicUpdate,
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = js.AddConsumer(DmUpdateStreamName, &nats.ConsumerConfig{
		Durable:        DmUpdateConsumeName,
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverSubject: nats.NewInbox(),
		//DeliverGroup:   ThingsDeliverGroup,
	})
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: js}, nil
}

func (n *NatsClient) TempModelUpdate(ctx context.Context, info *templateModel.TemplateInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = n.client.Publish(TopicUpdate, events.NewEventMsg(ctx, data))
	logx.WithContext(ctx).Infof("%s|info:%v,err:%v", utils.FuncName(),
		info, err)
	return err
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.Subscribe(TopicUpdate, func(msg *nats.Msg) {
		msg.Ack()
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Errorf("%v|GetEventMsg|subject:%v,data:%v",
				utils.FuncName(), msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		tempInfo := templateModel.TemplateInfo{}
		err := json.Unmarshal(emsg.GetData(), &tempInfo)
		if err != nil {
			logx.Errorf("%v|Unmarshal|subject:%v,data:%v",
				utils.FuncName(), msg.Subject, string(msg.Data))
			return
		}
		err = handle(ctx).TempModelClearCache(&tempInfo)
		logx.WithContext(ctx).Infof("%s|topic:%v,subject:%v,data:%v,err:%v", utils.FuncName(),
			TopicUpdate, msg.Subject, string(msg.Data), err)
	}, nats.Durable(DmUpdateConsumeName), nats.BindStream(DmUpdateStreamName))
	if err != nil {
		return err
	}
	return nil
}
