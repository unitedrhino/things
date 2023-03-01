package mqttTest

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type AssertInfo struct {
	SubTopic string
	PubTopic string
	Req      []byte
}

func Assert(mc mqtt.Client, info AssertInfo, assertFunc func(resp []byte) error) error {
	subChan := make(chan []byte)
	go mc.Subscribe(info.SubTopic, 1, func(client mqtt.Client, message mqtt.Message) {
		subChan <- message.Payload()
	}).WaitTimeout(3 * time.Second)
	time.Sleep(time.Second)
	err := mc.Publish(info.PubTopic, 1, false, info.Req).Error()
	if err != nil {
		return err
	}
	for {
		select {
		case respC := <-subChan:
			return assertFunc(respC)
		case <-time.After(5 * time.Second):
			return fmt.Errorf("subscribe time out")
		}
	}
}
