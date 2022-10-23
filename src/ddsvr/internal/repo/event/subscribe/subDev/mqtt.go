package subDev

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"strings"
	"time"
)

type (
	MqttClient struct {
		client mqtt.Client
	}
	//登录登出消息
	ConnectMsg struct {
		UserName string `json:"username"`
		Ts       int64  `json:"ts"`
		Address  string `json:"ipaddress"`
		ClientID string `json:"clientid"`
		Reason   string `json:"reason"`
	}
)

const (
	ActionConnected    = "connected"
	ActionDisconnected = "disconnected"
)
const (
	// ShareSubTopicPrefix emqx 共享订阅前缀 参考: https://docs.emqx.com/zh/enterprise/v4.4/advanced/shared-subscriptions.html
	ShareSubTopicPrefix = "$share/dd.rpc/"
	// TopicConnectStatus emqx 客户端上下线通知 参考: https://docs.emqx.com/zh/enterprise/v4.4/advanced/system-topic.html#客户端上下线事件
	TopicConnectStatus = ShareSubTopicPrefix + "$SYS/brokers/+/clients/#"

	TopicThing  = ShareSubTopicPrefix + devices.TopicHeadThing + "/#"
	TopicOta    = ShareSubTopicPrefix + devices.TopicHeadOta + "/#"
	TopicConfig = ShareSubTopicPrefix + devices.TopicHeadConfig + "/#"
	TopicSDKLog = ShareSubTopicPrefix + devices.TopicHeadLog + "/#"
	TopicShadow = ShareSubTopicPrefix + devices.TopicHeadShadow + "/#"
)

func newEmqClient(conf *conf.MqttConf) (SubDev, error) {
	mc, err := clients.NewMqttClient(conf)
	if err != nil {
		return nil, err
	}
	return &MqttClient{
		client: mc,
	}, nil
}

func (d *MqttClient) SubDevMsg(handle Handle) error {
	err := d.subscribeWithFunc(TopicConnectStatus, d.subscribeConnectStatus(handle))
	if err != nil {
		return err
	}
	err = d.subscribeWithFunc(TopicThing, func(ctx context.Context, topic string, payload []byte) error {
		return handle(ctx).Thing(topic, payload)
	})
	if err != nil {
		return err
	}
	err = d.subscribeWithFunc(TopicConfig, func(ctx context.Context, topic string, payload []byte) error {
		return handle(ctx).Config(topic, payload)
	})
	if err != nil {
		return err
	}
	err = d.subscribeWithFunc(TopicOta, func(ctx context.Context, topic string, payload []byte) error {
		return handle(ctx).Ota(topic, payload)
	})
	if err != nil {
		return err
	}
	err = d.subscribeWithFunc(TopicSDKLog, func(ctx context.Context, topic string, payload []byte) error {
		return handle(ctx).SDKLog(topic, payload)
	})
	if err != nil {
		return err
	}
	err = d.subscribeWithFunc(TopicShadow, func(ctx context.Context, topic string, payload []byte) error {
		return handle(ctx).Shadow(topic, payload)
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *MqttClient) subscribeConnectStatus(handle Handle) func(ctx context.Context, topic string, payload []byte) error {
	return func(ctx context.Context, topic string, payload []byte) error {
		var (
			msg ConnectMsg
			err error
		)
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			logx.Error(err)
			return err
		}
		do := devices.DevConn{
			UserName:  msg.UserName,
			Timestamp: msg.Ts, //毫秒时间戳
			Address:   msg.Address,
			ClientID:  msg.ClientID,
			Reason:    msg.Reason,
		}
		if strings.HasSuffix(topic, "/disconnected") {
			logx.WithContext(ctx).Infof("%s.disconnected topic:%v message:%v err:%v",
				utils.FuncName(), topic, string(payload), err)
			do.Action = ActionDisconnected
			err = handle(ctx).Disconnected(&do)
			if err != nil {
				logx.Error(err)
				return err
			}
		} else {
			do.Action = ActionConnected
			logx.WithContext(ctx).Infof("%s.connected topic:%v message:%v err:%v",
				utils.FuncName(), topic, string(payload), err)
			err = handle(ctx).Connected(&do)
			if err != nil {
				logx.Error(err)
				return err
			}
		}
		return nil
	}

}

func (d *MqttClient) subscribeWithFunc(topic string, handle func(ctx context.Context, topic string, payload []byte) error) error {
	return d.client.Subscribe(topic,
		1, func(client mqtt.Client, message mqtt.Message) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			//ddsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
			//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
			ctx, span := traces.StartSpan(ctx, message.Topic(), "")
			defer span.End()
			startTime := timex.Now()
			duration := timex.Since(startTime)
			err := handle(ctx, message.Topic(), message.Payload())
			logx.WithContext(ctx).WithDuration(duration).Infof(
				"subscribeWithFunc.Subscribe.publish topic:%v message:%v err:%v",
				message.Topic(), string(message.Payload()), err)
		}).Error()
}
