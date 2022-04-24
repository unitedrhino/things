package dataUpdate

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/i-Things/things/src/dmsvr/dmDef"
	"github.com/i-Things/things/src/dmsvr/internal/domain/templateModel"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client nats.JetStreamContext
	}
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
		Name: dmDef.DmUpdateStreamName,
		Subjects: []string{
			dmDef.TopicUpdate,
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = js.AddConsumer(dmDef.DmUpdateStreamName, &nats.ConsumerConfig{
		Durable:   dmDef.DmUpdateConsumeName,
		AckPolicy: nats.AckExplicitPolicy,
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
	_, err = n.client.Publish(dmDef.TopicUpdate, events.NewEventMsg(ctx, data))
	return err
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.Subscribe(dmDef.DmUpdateStreamName, func(msg *nats.Msg) {
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
			ddExport.TopicDevPublishAll, msg.Subject, string(msg.Data), err)
	})
	if err != nil {
		return err
	}
	return nil
}
