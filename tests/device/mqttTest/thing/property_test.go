package thing

import (
	"fmt"
	"gitee.com/unitedrhino/things/tests/device/mqttTest"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
		c := &mqttTest.DefaultClientInfo
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

func TestProperty(t *testing.T) {
	t.Logf("TestProperty")
	_, mc := GetInit(t)
	subTopic := fmt.Sprintf("$thing/down/property/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("$thing/up/property/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	mc.Subscribe(subTopic, 1, func(client mqtt.Client, message mqtt.Message) {
		t.Log(message)
	})
	pubPayload := `
{                     
    "method": "report",            
    "msgToken": "afwegafeegfa",   
    "params": 	{"GPS_Info":{"longtitude":0.125,"latitude":0.005}}
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
