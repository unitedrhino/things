package innerLink

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
)

type Direct struct {
}

func NewDirect() {

}

func (d Direct) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) Subscribe(handle Handle) error {
	//TODO implement me
	panic("implement me")
}

func (d Direct) ReqToDeviceSync(ctx context.Context, reqTopic, respTopic string, req *deviceSend.DeviceReq, productID, deviceName string) (*deviceSend.DeviceResp, error) {
	//TODO implement me
	panic("implement me")
}
