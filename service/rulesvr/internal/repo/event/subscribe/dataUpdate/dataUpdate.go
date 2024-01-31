package dataUpdate

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/events"
)

type (
	DataUpdate interface {
		Subscribe(handle Handle) error
	}
	Handle       func(ctx context.Context) UpdateHandle
	UpdateHandle interface {
		ProductSchemaUpdate(info *events.DeviceUpdateInfo) error
		SceneInfoDelete(info *events.ChangeInfo) error
		SceneInfoUpdate(info *events.ChangeInfo) error
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
