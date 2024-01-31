package thing

import (
	"encoding/base64"
	"fmt"
	"gitee.com/i-Things/core/shared/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/tests/device/mqttTest"
	"github.com/spf13/cast"
	"sync"
	"testing"
	"time"
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

func TestGatewayOperationRegister(t *testing.T) {
	clientInfo, mc := GetInit(t)
	subTopic := fmt.Sprintf("$gateway/down/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubTopic := fmt.Sprintf("$gateway/up/operation/%s/%s", clientInfo.ProductID, clientInfo.DeviceName)
	pubPayload := `
{
  "method": "register",
  "payload": {
    "devices": [
      {
        "productID": "255fCKEJ02I",
        "deviceName": "subdeviceaa5"
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
	subDevice := mqttTest.SubDevice
	random := utils.Random(5, 0)
	ts := time.Now().Unix()
	sign := fmt.Sprintf("%v;%v;%v;%v", subDevice.ProductID, subDevice.DeviceName, random, ts)
	pwd, _ := base64.StdEncoding.DecodeString(subDevice.DeviceSecret)
	signature := utils.HmacSha256(sign, pwd)
	pubPayload := `
	{
  "method": "bind",
  "payload": {
    "devices": [
      {
        "productID": "%s",
        "deviceName": "%s",
        "signature": "%s",
        "random": %d,
        "timestamp": %v,
        "signMethod": "hmacsha256"
      }
    ]
  }
}
`
	pub := fmt.Sprintf(pubPayload, subDevice.ProductID, subDevice.DeviceName, signature, cast.ToInt(random), ts)
	err := mqttTest.Assert(mc, mqttTest.AssertInfo{
		SubTopic: subTopic,
		PubTopic: pubTopic,
		Req:      []byte(pub),
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
