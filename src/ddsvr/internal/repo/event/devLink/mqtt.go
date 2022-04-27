package devLink

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
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
	ActionLogin  = "onLogin"
	ActionLogout = "onLogout"
)

func NewEmqClient(conf *conf.MqttConf) (DevLink, error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range conf.Brokers {
		opts.AddBroker(broker)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		logx.Info("GenerateUUID failure")
		return nil, err
	}
	opts.SetClientID(conf.ClientID + "/" + uuid).SetUsername(conf.User).
		SetPassword(conf.Pass).SetAutoReconnect(true).SetConnectRetry(true)
	opts.OnConnect = func(client mqtt.Client) {
		logx.Info("Connected")
	}
	mc := mqtt.NewClient(opts)
	er := mc.Connect().WaitTimeout(5 * time.Second)
	if er == false {
		logx.Info("Connect failure")
		return nil, fmt.Errorf("mqtt client connect failure")
	}

	c := &MqttClient{
		client: mc,
	}

	return c, nil
}

func (d *MqttClient) SubScribe(handle Handle) error {
	err := d.client.Subscribe("$share/dd.rpc/$SYS/brokers/+/clients/#",
		1, func(client mqtt.Client, message mqtt.Message) {
			var (
				msg ConnectMsg
				err error
			)
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			err = json.Unmarshal(message.Payload(), &msg)
			if err != nil {
				logx.Error(err)
				return
			}
			do := devices.DevConn{
				UserName:  msg.UserName,
				Timestamp: msg.Ts, //毫秒时间戳
				Address:   msg.Address,
				ClientID:  msg.ClientID,
				Reason:    msg.Reason,
			}
			if strings.HasSuffix(message.Topic(), "/disconnected") {
				logx.WithContext(ctx).Infof("%s|disconnected|topic:%v,message:%v,err:%v",
					utils.FuncName(), message.Topic(), string(message.Payload()), err)
				do.Action = ActionLogout
				err = handle(ctx).Disconnected(&do)
				if err != nil {
					logx.Error(err)
				}
			} else {
				do.Action = ActionLogin
				logx.WithContext(ctx).Infof("%s|connected|topic:%v,message:%v,err:%v",
					utils.FuncName(), message.Topic(), string(message.Payload()), err)
				err = handle(ctx).Connected(&do)
				if err != nil {
					logx.Error(err)
				}
			}
		}).Error()
	if err != nil {
		return err
	}
	err = d.client.Subscribe("$share/dd.rpc/$thing/#",
		1, func(client mqtt.Client, message mqtt.Message) {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			err := handle(ctx).Publish(message.Topic(), message.Payload())
			logx.WithContext(ctx).Infof("%s|publish|topic:%v,message:%v,err:%v",
				utils.FuncName(), message.Topic(), string(message.Payload()), err)

		}).Error()
	return err
}

func (d *MqttClient) Publish(ctx context.Context, topic string, payload []byte) error {
	return d.client.Publish(topic, 1, false, payload).Error()
}
