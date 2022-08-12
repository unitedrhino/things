package pubInner

import (
	//"github.com/i-Things/things/src/ddsvr/internal/domain"
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
)

type ConnType int8

const (
	Connect ConnType = iota
	DisConnect
)

type (
	PubInner interface {
		// DevPubThing 向内部发布设备发布的物模型消息
		DevPubThing(ctx context.Context, publishMsg *devices.DevPublish) error
		// DevPubOta 发布ota升级相关消息
		DevPubOta(ctx context.Context, publishMsg *devices.DevPublish) error
		// DevPubShadow 发布设备影子消息
		DevPubShadow(ctx context.Context, publishMsg *devices.DevPublish) error
		// DevPubSDKLog 发布设备调试日志
		DevPubSDKLog(ctx context.Context, publishMsg *devices.DevPublish) error
		// DevPubConfig 发布设备配置相关消息
		DevPubConfig(ctx context.Context, publishMsg *devices.DevPublish) error
		// PubConn 向内部发布连接及断连消息
		PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error
	}
)

func NewPubInner(c conf.EventConf) (PubInner, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
