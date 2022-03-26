package innerLink

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/i-Things/things/src/dmsvr/dmDef"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceSend"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
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

func (n *NatsClient) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	return n.client.Publish(ddExport.TopicInnerPublish, ddExport.PublishToDev(ctx, topic, payload))
}
func (n *NatsClient) Subscribe(handle Handle) error {
	n.client.QueueSubscribe(ddExport.TopicDevPublishAll, dmDef.SvrName, func(msg *nats.Msg) {
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ele, err := n.getDevPublish(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Publish(ele)
		logx.WithContext(ctx).Info(ddExport.TopicDevPublishAll, msg.Subject, string(msg.Data), err)
	})
	n.client.QueueSubscribe(ddExport.TopicDevConnected, dmDef.SvrName, func(msg *nats.Msg) {
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ele, err := n.getDevConnMsg(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Connected(ele)
		logx.WithContext(ctx).Info(ddExport.TopicDevConnected, msg.Subject, string(msg.Data), err)
	})
	n.client.QueueSubscribe(ddExport.TopicDevDisconnected, dmDef.SvrName, func(msg *nats.Msg) {
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ele, err := n.getDevConnMsg(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Disconnected(ele)
		logx.WithContext(ctx).Info(ddExport.TopicDevDisconnected, msg.Subject, string(msg.Data), err)
	})
	return nil
}
func (n *NatsClient) getDevPublish(ctx context.Context, data []byte) (*deviceSend.Elements, error) {
	pubInfo := ddExport.DevPublish{}
	err := json.Unmarshal(data, &pubInfo)
	if err != nil {
		logx.WithContext(ctx).Error("getDevPublish", string(data), err)
		return nil, err
	}
	ele := deviceSend.Elements{
		Topic:      pubInfo.Topic,
		Payload:    pubInfo.Payload,
		Timestamp:  pubInfo.Timestamp,
		ProductID:  pubInfo.ProductID,
		DeviceName: pubInfo.DeviceName,
	}
	return &ele, nil
}

func (n *NatsClient) getDevConnMsg(ctx context.Context, data []byte) (*deviceSend.Elements, error) {
	logInfo := ddExport.DevConn{}
	err := json.Unmarshal(data, &logInfo)
	if err != nil {
		logx.WithContext(ctx).Error("getDevConnMsg", string(data), err)
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
