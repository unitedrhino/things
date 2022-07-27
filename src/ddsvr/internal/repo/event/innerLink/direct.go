package innerLink

import (
	"context"
	"github.com/i-Things/things/shared/devices"
)

type Direct struct {
}

func NewDirect() (InnerLink, error) {
	return &Direct{}, nil
}

func (d Direct) DevPubThing(ctx context.Context, publishMsg *devices.DevPublish) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) DevPubOta(ctx context.Context, publishMsg *devices.DevPublish) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) DevPubShadow(ctx context.Context, publishMsg *devices.DevPublish) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) DevPubSDKLog(ctx context.Context, publishMsg *devices.DevPublish) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) DevPubConfig(ctx context.Context, publishMsg *devices.DevPublish) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) Subscribe(handle Handle) error {
	//TODO implement me
	panic("implement me")
}
