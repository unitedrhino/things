package subInner

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
)

type (
	SubInner interface {
		SubToDevMsg(handle Handle) error
	}
	Handle         func(ctx context.Context) InnerSubHandle
	InnerSubHandle interface {
		PublishToDev(info *devices.InnerPublish) error
	}
)

func NewSubInner(c conf.EventConf) (SubInner, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
