package clients

import (
	"crypto/tls"
	"fmt"
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
	mqttClient   mqtt.Client
	// MqttSetOnConnectHandler 如果会话断开可以通过该回调函数来重新订阅消息
	//不使用mqtt的clean session是因为会话保持期间共享订阅也会给离线的客户端,这会导致在线的客户端丢失消息
	MqttSetOnConnectHandler func()
)

func NewMqttClient(conf *conf.MqttConf) (mc mqtt.Client, err error) {
	mqttInitOnce.Do(func() {
		for i := 3; i > 0; i-- {
			mc, err = initMqtt(conf)
			if err != nil { //出现并发情况的时候可能iThings的http还没启动完毕
				logx.Errorf("mqtt 连接失败 重试剩余次数:%v", i-1)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		if err != nil {
			logx.Errorf("mqtt 连接失败 conf:%#v  err:%v", conf, err)
			os.Exit(-1)
		}
		mqttClient = mc
	})
	return mqttClient, err
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
	opts.SetClientID(conf.ClientID + "/" + uuid).SetUsername(conf.User).
		SetPassword(conf.Pass).SetAutoReconnect(true).SetConnectRetry(true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Info("mqtt client Connected")
		if MqttSetOnConnectHandler != nil {
			MqttSetOnConnectHandler()
		}
	})
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
