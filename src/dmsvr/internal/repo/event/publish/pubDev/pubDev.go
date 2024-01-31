package pubDev

import (
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
)

type (
	CompareMsg func(payload []byte) bool

	PubDev interface {
		PublishToDev(ctx context.Context, msg *deviceMsg.PublishMsg) error
		ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, compareMsg CompareMsg) ([]byte, error)
	}
)

func NewPubDev(c conf.EventConf) (PubDev, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	case conf.EventModeNatsJs:
		return newNatsJsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
