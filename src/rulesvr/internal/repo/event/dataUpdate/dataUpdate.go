package dataUpdate

import (
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/errors"
)

type (
	DataUpdate interface {
		UpdateWithTopic(ctx context.Context, topic string, info any) error
	}
)

func NewDataUpdate(c conf.EventConf) (DataUpdate, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	case conf.EventModeNatsJs:
		return newNatsJsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)

}
