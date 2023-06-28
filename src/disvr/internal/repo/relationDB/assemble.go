package relationDB

import (
	"github.com/i-Things/things/src/disvr/internal/domain/shadow"
)

func ToShadowPo(info *shadow.Info) *DiDeviceShadow {
	return &DiDeviceShadow{
		ID:                info.ID,
		ProductID:         info.ProductID,
		DeviceName:        info.DeviceName,
		DataID:            info.DataID,
		UpdatedDeviceTime: info.UpdatedDeviceTime,
		Value:             info.Value,
	}
}
func ToShadowDo(in *DiDeviceShadow) *shadow.Info {
	return &shadow.Info{
		ID:                in.ID,
		ProductID:         in.ProductID,
		DeviceName:        in.DeviceName,
		DataID:            in.DataID,
		Value:             in.Value,
		UpdatedDeviceTime: in.UpdatedDeviceTime,
		CreatedTime:       in.CreatedTime,
		UpdatedTime:       in.UpdatedTime,
	}
}
func ToShadowsDo(in []*DiDeviceShadow) []*shadow.Info {
	if in == nil {
		return nil
	}
	var ret []*shadow.Info
	for _, v := range in {
		ret = append(ret, ToShadowDo(v))
	}
	return ret
}
