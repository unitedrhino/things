package pubDev

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
)

type (
	// PubDev 发送消息到设备
	PubDev interface {
		Publish(ctx context.Context, topic string, payload []byte) error
		CheckIsOnline(ctx context.Context, clientID string) (bool, error)
	}
)

func Check(conf conf.DevLinkConf) error {
	if conf.Mqtt == nil {
		return fmt.Errorf("DevLinkConf need")
	}
	return nil
}

func NewPubDev(conf conf.DevLinkConf) (PubDev, error) {
	if err := Check(conf); err != nil {
		return nil, err
	}
	return newEmqClient(conf.Mqtt)
}
