package mqttTest

import (
	"testing"
)

func TestGetMqttClient(t *testing.T) {
	_, err := GetMqttClient(&ClientInfo{
		MqttBrokers:  []string{"tcp://106.15.225.172:1883"},
		ProductID:    "254pwnKQsvK",
		DeviceName:   "test5",
		DeviceSecret: "6skuocmEga94+OhVYRGWUphWlX0=",
	})
	if err != nil {
		t.Error(err)
	}
}
