package collect

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToUserCollectDeviceSavePb(in *types.UserCollectDeviceSave) *ud.UserCollectDeviceSave {
	if in == nil {
		return nil
	}
	var devices []*ud.DeviceCore
	for _, v := range in.Devices {
		devices = append(devices, &ud.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return &ud.UserCollectDeviceSave{
		Devices: devices,
	}
}
