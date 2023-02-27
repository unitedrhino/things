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
		c := &mqttTest.DefaultGateway
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

func TestGatewayStatusOnline(t *testing.T) {

	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("$gateway/down/status/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("$gateway/up/status/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
	{
  "method": "online",
  "payload": {
    "devices": [
       {
            "productID": "255fCKEJ02I",
            "deviceName": "test2"
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

func TestGatewayOperationDescribeSubDevices(t *testing.T) {

	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("$gateway/down/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("$gateway/up/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
	{
	  "method": "describeSubDevices"
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

func TestGatewayOperationBind(t *testing.T) {

	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("$gateway/down/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("$gateway/up/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
	{
	  "method": "bind",
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

func TestGatewayOperationUnBind(t *testing.T) {

	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("$gateway/down/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("$gateway/up/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
	{
	  "method": "unbind",
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

func TestGatewayOperationAccessSubDev(t *testing.T) {
	_, mc := GetInit(t)
	var (
		subDevProductID  = "255fCKEJ02I"
		subDevDeviceName = "test2"
	)

	subTopic := fmt.Sprintf("$thing/down/property/%s/%s", subDevProductID, subDevDeviceName)
	pubTopic := fmt.Sprintf("$thing/up/property/%s/%s", subDevProductID, subDevDeviceName)
	pubPayload := `
	{                     
    "method": "report",            
    "clientToken": "afwegafeegfa",   
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
