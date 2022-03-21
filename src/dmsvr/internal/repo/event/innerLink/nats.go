package innerLink

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/ddsvr/ddDef"
	"github.com/i-Things/things/src/dmsvr/dmDef"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceSend"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

func NewNatsClient(conf conf.NatsConf) (InnerLink, error) {
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

func (n *NatsClient) Publish(ctx context.Context, topic string, payload []byte) error {
	return n.client.Publish(topic, payload)
}
func (n *NatsClient) Subscribe(handle Handle) error {
	n.client.QueueSubscribe(ddDef.TopicDevPublish, dmDef.SvrName, func(msg *nats.Msg) {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		ele, err := n.getDevLogInOut(ctx, msg)
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Publish(ele)
		logx.WithContext(ctx).Info(ddDef.TopicDevPublish, msg.Subject, string(msg.Data), err)
	})
	n.client.QueueSubscribe(ddDef.TopicDevLogin, dmDef.SvrName, func(msg *nats.Msg) {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		ele, err := n.getDevLogInOut(ctx, msg)
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Login(ele)
		logx.WithContext(ctx).Info(ddDef.TopicDevLogin, msg.Subject, string(msg.Data), err)
	})
	n.client.QueueSubscribe(ddDef.TopicDevLogout, dmDef.SvrName, func(msg *nats.Msg) {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		ele, err := n.getDevLogInOut(ctx, msg)
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Logout(ele)
		logx.WithContext(ctx).Info(ddDef.TopicDevLogout, msg.Subject, string(msg.Data), err)
	})
	return nil
}
func (n *NatsClient) getDevPublish(ctx context.Context, msg *nats.Msg) (*deviceSend.Elements, error) {
	data := def.MsgHead{}
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return nil, err
	}
	pubInfo := ddDef.DevPublish{}
	err = json.Unmarshal(data.Data, &pubInfo)
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return nil, err
	}
	ele := deviceSend.Elements{
		Topic:     pubInfo.Topic,
		Payload:   pubInfo.Payload,
		Timestamp: pubInfo.Timestamp,
	}
	return &ele, nil
}

func (n *NatsClient) getDevLogInOut(ctx context.Context, msg *nats.Msg) (*deviceSend.Elements, error) {
	data := def.MsgHead{}
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return nil, err
	}
	logInfo := ddDef.DevLogInOut{}
	err = json.Unmarshal(data.Data, &logInfo)
	if err != nil {
		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
		return nil, err
	}
	ele := deviceSend.Elements{
		ClientID:  logInfo.ClientID,
		Username:  logInfo.UserName,
		Timestamp: logInfo.Timestamp,
		Address:   logInfo.Address,
		Action:    logInfo.Action,
		Reason:    logInfo.Reason,
	}
	return &ele, nil
}
