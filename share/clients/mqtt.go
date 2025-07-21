package clients

import (
	"crypto/tls"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	cfg     *conf.MqttConf
}

func NewMqttClient(conf *conf.MqttConf) (mcs *MqttClient, err error) {
	mqttInitOnce.Do(func() {
		var clients []mqtt.Client
		var start = time.Now()
		for len(clients) < conf.ConnNum {
			var (
				mc mqtt.Client
			)
			var tryTime = 5
			for i := tryTime; i > 0; i-- {
				mc, err = initMqtt(conf)
				logx.Infof("mqtt_client initMqtt2 mc:%v err:%v", mc, err)
				if err != nil { //出现并发情况的时候可能联犀的http还没启动完毕
					logx.Errorf("mqtt_client 连接失败 重试剩余次数:%v", i-1)
					time.Sleep(time.Second * time.Duration(tryTime) / time.Duration(i))
					continue
				}
				break
			}
			if err != nil {
				logx.Errorf("mqtt_client 连接失败 conf:%#v  err:%v", conf, err)
				os.Exit(-1)
			}
			clients = append(clients, mc)
			var cli = MqttClient{clients: clients, cfg: conf}
			mqttClient = &cli
			logx.Infof("mqtt_client 连接完成 clientNum:%v use:%s", len(clients), time.Now().Sub(start))
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
	logx.Infof("mqtt_client_subscribe clientNum:%v topic:%v", len(clients), topic)
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
	uuid := uuid.NewString()
	clientID := conf.ClientID + "_" + uuid
	logx.Infof("mqtt_client initMqtt conf:%#v clientID:%v brokers:%#v stack=%s", conf, clientID, opts.Servers, utils.Stack(1, 10))
	opts.SetClientID(clientID).SetUsername(conf.User).SetPassword(conf.Pass)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Infof("mqtt_client Connected clientID:%v", clientID)
		if mqttSetOnConnectHandler != nil {
			mqttSetOnConnectHandler(client)
		}
	})
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		logx.Infof("mqtt_client Reconnecting clientID:%#v", options)
		if mqttSetOnConnectHandler != nil {
			mqttSetOnConnectHandler(client)
		}
	})

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(30 * time.Second) //意外离线的重连参数
	opts.SetConnectRetry(true).SetConnectRetryInterval(5 * time.Second)   //首次连接的重连参数

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt_client 正在尝试连接 broker:%v clientID:%v", utils.Fmt(broker), clientID)
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt_client 连接丢失 err:%v  clientID:%v", utils.Fmt(err), clientID)
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		logx.Errorf("mqtt_client 连接失败超时")
		err = fmt.Errorf("mqtt_client 连接失败")
		return nil, err
	}
	return
}
