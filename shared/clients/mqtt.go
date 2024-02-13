package clients

import (
	"crypto/tls"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	mqttInitOnce sync.Once
	mqttClient   *MqttClient
	// mqttSetOnConnectHandler 如果会话断开可以通过该回调函数来重新订阅消息
	//不使用mqtt的clean session是因为会话保持期间共享订阅也会给离线的客户端,这会导致在线的客户端丢失消息
	mqttSetOnConnectHandler func(cli mqtt.Client)
)

type MqttClient struct {
	clients []mqtt.Client
}

func NewMqttClient(conf *conf.MqttConf) (mcs *MqttClient, err error) {
	mqttInitOnce.Do(func() {
		var clients []mqtt.Client
		for len(clients) < conf.ConnNum {
			var (
				mc mqtt.Client
			)
			var tryTime = 5
			for i := tryTime; i > 0; i-- {
				mc, err = initMqtt(conf)
				if err != nil { //出现并发情况的时候可能iThings的http还没启动完毕
					logx.Errorf("mqtt 连接失败 重试剩余次数:%v", i-1)
					time.Sleep(time.Second * time.Duration(tryTime) / time.Duration(i))
					continue
				}
				break
			}
			if err != nil {
				logx.Errorf("mqtt 连接失败 conf:%#v  err:%v", conf, err)
				os.Exit(-1)
			}
			clients = append(clients, mc)
			var cli = MqttClient{clients: clients}
			mqttClient = &cli
		}
	})
	return mqttClient, err
}

func SetMqttSetOnConnectHandler(f func(cli mqtt.Client)) {
	mqttSetOnConnectHandler = f
}

func (m MqttClient) Subscribe(cli mqtt.Client, topic string, qos byte, callback mqtt.MessageHandler) error {
	var clients = m.clients
	if cli != nil {
		clients = []mqtt.Client{cli}
	}
	for _, c := range clients {
		err := c.Subscribe(topic, qos, callback).Error()
		if err != nil {
			return errors.System.AddDetail(err)
		}
	}
	return nil
}

func (m MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	id := rand.Intn(len(m.clients))
	return m.clients[id].Publish(topic, qos, retained, payload).Error()
}

func initMqtt(conf *conf.MqttConf) (mc mqtt.Client, err error) {
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
	opts.SetClientID(conf.ClientID + "/" + uuid).SetUsername(conf.User).SetPassword(conf.Pass)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Info("mqtt client Connected")
		if mqttSetOnConnectHandler != nil {
			mqttSetOnConnectHandler(client)
		}
	})

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(5 * time.Second) //意外离线的重连参数
	opts.SetConnectRetry(true).SetConnectRetryInterval(1 * time.Second)  //首次连接的重连参数

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt 正在尝试连接 broker:%v", utils.Fmt(broker))
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt 连接丢失 err:%v", utils.Fmt(err))
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		logx.Info("mqtt 连接失败")
		err = fmt.Errorf("mqtt 连接失败")
		return
	}
	return
}
