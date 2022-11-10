package pubDev

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
)

type (
	PubDev interface {
		PublishToDev(ctx context.Context, topic string, payload []byte) error
		ReqToDeviceSync(ctx context.Context, reqTopic, respTopic string, req *msgThing.Req,
			productID, deviceName string) (*msgThing.Resp, error)
	}
)

func NewPubDev(c conf.EventConf) (PubDev, error) {
	if c.Mode == conf.EventModeNats {
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
