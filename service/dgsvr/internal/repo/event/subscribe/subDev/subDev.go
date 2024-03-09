package subDev

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/devices"
)

type (
	SubDev interface {
		SubDevMsg(handle Handle) error
	}
	Handle       func(ctx context.Context) DevSubHandle
	DevSubHandle interface {
		ExtMsg(topic string, payload []byte) error
		Msg(topic string, payload []byte) error
		Connected(out *devices.DevConn) error
		Disconnected(out *devices.DevConn) error
	}
)

func Check(conf conf.DevLinkConf) error {
	if conf.Mqtt == nil {
		return fmt.Errorf("DevLinkConf need")
	}
	return nil
}

func NewSubDev(conf conf.DevLinkConf) (SubDev, error) {
	if err := Check(conf); err != nil {
		return nil, err
	}
	return newEmqClient(conf.Mqtt)
}
