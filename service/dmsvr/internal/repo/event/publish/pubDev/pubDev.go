package pubDev

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/errors"
)

type (
	CompareMsg func(payload []byte) bool

	PubDev interface {
		PublishToDev(ctx context.Context, msg *deviceMsg.PublishMsg) error
		ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, compareMsg CompareMsg) ([]byte, error)
	}
)

func NewPubDev(c conf.EventConf, nodeID int64) (PubDev, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c, nodeID)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
