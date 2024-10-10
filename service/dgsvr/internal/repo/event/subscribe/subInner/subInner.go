package subInner

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
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

func NewSubInner(c conf.EventConf, nodeID int64) (SubInner, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c, nodeID)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
