package protocol

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"time"
)

type DevHandle func(ctx context.Context, topic string, payload []byte) error

type MqttProtocol struct {
	*LightProtocol
	client       *clients.MqttClient
	DevSubHandle map[string]DevHandle
}

func NewMqttProtocol(c conf.EventConf, pi *dm.ProtocolInfo, pc *LightProtocolConf, mqttc conf.DevLinkConf) (*MqttProtocol, error) {
	if mqttc.Mqtt == nil {
		return nil, errors.Parameter.AddMsg("DevLinkConf need")
	}
	mc, err := clients.NewMqttClient(mqttc.Mqtt)
	if err != nil {
		return nil, err
	}
	lightProto, err := NewLightProtocol(c, pi, pc)
	if err != nil {
		return nil, err
	}
	return &MqttProtocol{LightProtocol: lightProto, client: mc, DevSubHandle: make(map[string]DevHandle)}, nil
}

func (m *MqttProtocol) SubscribeDevMsg(topic string, handle DevHandle) error {
	newTopic := fmt.Sprintf("$share/%s/%s", m.ServerName, topic)
	m.DevSubHandle[newTopic] = handle
	return nil
}

func (m *MqttProtocol) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	err := m.client.Publish(topic, 1, false, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Publish failure err:%v topic:%v", err, topic)
	}
	return err
}

func (m *MqttProtocol) Start() error {
	err := m.subscribes(nil)
	if err != nil {
		return err
	}
	clients.SetMqttSetOnConnectHandler(func(cli mqtt.Client) {
		err := m.subscribes(cli)
		if err != nil {
			logx.Errorf("%s.mqttSetOnConnectHandler.subscribes err:%v", utils.FuncName(), err)
		}
	})
	return m.LightProtocol.Start()
}

func (m *MqttProtocol) subscribes(cli mqtt.Client) error {
	for topic, handle := range m.DevSubHandle {
		err := m.subscribeWithFunc(cli, topic, handle)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MqttProtocol) subscribeWithFunc(cli mqtt.Client, topic string, handle func(ctx context.Context, topic string, payload []byte) error) error {
	return m.client.Subscribe(cli, topic,
		1, func(client mqtt.Client, message mqtt.Message) {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
				defer cancel()
				utils.Recover(ctx)
				//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
				//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
				ctx, span := ctxs.StartSpan(ctx, message.Topic(), "")
				defer span.End()
				startTime := timex.Now()
				duration := timex.Since(startTime)
				ctx = ctxs.WithRoot(ctx)
				err := handle(ctx, message.Topic(), message.Payload())
				if err != nil {
					logx.WithContext(ctx).Errorf("%s.handle failure err:%v topic:%v", utils.FuncName(), err, topic)
				}
				logx.WithContext(ctx).WithDuration(duration).Infof(
					"subscribeWithFunc.Subscribe.publish topic:%v message:%v err:%v",
					message.Topic(), string(message.Payload()), err)
			}()

		})
}
