package clients

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

var (
	mqttInitOnce sync.Once
	mqttClient   mqtt.Client
)

func NewMqttClient(conf *conf.MqttConf) (mc mqtt.Client, err error) {
	mqttInitOnce.Do(func() {
		opts := mqtt.NewClientOptions()
		for _, broker := range conf.Brokers {
			opts.AddBroker(broker)
		}
		uuid, er := uuid.GenerateUUID()
		if er != nil {
			logx.Info("GenerateUUID failure")
			err = er
			return
		}
		opts.SetClientID(conf.ClientID + "/" + uuid).SetUsername(conf.User).
			SetPassword(conf.Pass).SetAutoReconnect(true).SetConnectRetry(true)
		opts.OnConnect = func(client mqtt.Client) {
			logx.Info("mqtt client Connected")
		}
		mc := mqtt.NewClient(opts)
		er2 := mc.Connect().WaitTimeout(5 * time.Second)
		if er2 == false {
			logx.Info("mqtt Connect failure")
			err = fmt.Errorf("mqtt client connect failure")
			return
		}
		mqttClient = mc
	})
	return mqttClient, err
}
