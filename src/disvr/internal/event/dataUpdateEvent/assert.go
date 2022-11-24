package dataUpdateEvent

import (
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgGateway"
)

func ToGatewayPayload(status int32, in []*events.DeviceCore) *msgGateway.GatewayPayload {
	var ret = msgGateway.GatewayPayload{Status: status}
	for _, v := range in {
		ret.Devices = append(ret.Devices, &msgGateway.Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return &ret
}
