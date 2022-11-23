package thing

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/test/device/mqttTest"
	"sync"
	"testing"
)

var (
	getOnce    = sync.Once{}
	clientInfo *mqttTest.ClientInfo
	mc         mqtt.Client
)

func GetInit(t *testing.T) (*mqttTest.ClientInfo, mqtt.Client) {
	getOnce.Do(func() {
		c := &mqttTest.DefaultRemoteConfig
		mcc, err := mqttTest.GetMqttClient(c)
		if err != nil {
			t.Error(err)
			return
		}
		mc = mcc
		clientInfo = c
	})

	return clientInfo, mc
}

func TestRemoteConfigPush(t *testing.T) {

	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("config/down/push/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("config/up/push/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
	{
	  "method": "push",
 "payload": {
    "devices": [
      {
            "productID": "255fCKEJ02I",
            "deviceName": "test1"
        }
    ]
  }
	}
`
	err := mqttTest.Assert(mc, mqttTest.AssertInfo{
		SubTopic: subTopic,
		PubTopic: pubTopic,
		Req:      []byte(pubPayload),
	}, func(resp []byte) error {
		t.Log(string(resp))
		return nil
	})
	t.Log(err)
}

func TestRemoteConfigGet(t *testing.T) {

	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("config/down/get/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("config/up/get/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
	{
	  "method": "get",
 "payload": {
    "devices": [
      {
            "productID": "255fCKEJ02I",
            "deviceName": "test1"
        }
    ]
  }
	}
`
	err := mqttTest.Assert(mc, mqttTest.AssertInfo{
		SubTopic: subTopic,
		PubTopic: pubTopic,
		Req:      []byte(pubPayload),
	}, func(resp []byte) error {
		t.Log(string(resp))
		return nil
	})
	t.Log(err)
}
