package pubDev

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	deviceSend "github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
)

type (
	PubDev interface {
		PublishToDev(ctx context.Context, topic string, payload []byte) error
		ReqToDeviceSync(ctx context.Context, reqTopic, respTopic string, req *deviceSend.DeviceReq,
			productID, deviceName string) (*deviceSend.DeviceResp, error)
	}
)

func NewPubDev(c conf.EventConf) (PubDev, error) {
	if c.Mode == conf.EventModeNats {
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
