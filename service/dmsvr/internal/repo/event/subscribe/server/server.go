package server

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/share/errors"
)

type (
	Server interface {
		Subscribe(handle Handle) error
	}
	Handle       func(ctx context.Context) ServerHandle
	ServerHandle interface {
		ActionCheck(req *deviceMsg.PublishMsg) error
	}
)

func NewServer(c conf.EventConf, nodeID int64) (Server, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c, nodeID)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
