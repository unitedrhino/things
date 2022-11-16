package mqttTest

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"net/url"
	"time"
)

type ClientInfo struct {
	MqttBrokers  []string //mqtt服务器节点
	ProductID    string
	DeviceName   string
	DeviceSecret string //设备秘钥
}

var (
	DefaultClientInfo = ClientInfo{
		MqttBrokers:  []string{"tcp://106.15.225.172:1883"},
		ProductID:    "254pwnKQsvK",
		DeviceName:   "test5",
		DeviceSecret: "6skuocmEga94+OhVYRGWUphWlX0=",
	}
	DefaultGateway = ClientInfo{
		MqttBrokers:  []string{"tcp://106.15.225.172:1883"},
		ProductID:    "255fCsZtKEM",
		DeviceName:   "test1",
		DeviceSecret: "vsjl+0R8/kpkLd0PJ0my1HN8XDg=",
	}
)

func GetMqttClient(c *ClientInfo) (mc mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range c.MqttBrokers {
		opts.AddBroker(broker)
	}
	clientID, userName, pwd := deviceAuth.GenSecretDeviceInfo(deviceAuth.HmacSha256, c.ProductID, c.DeviceName, c.DeviceSecret)
	opts.SetClientID(clientID).SetUsername(userName).
		SetPassword(pwd).SetAutoReconnect(true).SetConnectRetry(true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Info("mqtt client Connected")
	})
	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt client connect attempt broker:%v", utils.Fmt(broker))
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt client connection lost err:%v", utils.Fmt(err))
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		logx.Info("mqtt Connect failure")
		err = fmt.Errorf("mqtt client connect failure")
		return
	}
	return mc, err
}
