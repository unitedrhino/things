package share

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToSharePb(in *types.UserDeviceShareInfo) *dm.UserDeviceShareInfo {
	if in == nil {
		return nil
	}
	return &dm.UserDeviceShareInfo{
		Id: in.ID,
		Device: &dm.DeviceCore{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
		SharedUserID: in.SharedUserID,
		ProjectID:    in.ProjectID,
		NormalPerm:   in.NormalPerm,
		SystemPerm:   in.SystemPerm,
	}
}

func ToShareTypes(in *dm.UserDeviceShareInfo) *types.UserDeviceShareInfo {
	if in == nil {
		return nil
	}
	return &types.UserDeviceShareInfo{
		ID: in.Id,
		Device: types.DeviceCore{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
		CreatedTime:       in.CreatedTime,
		SharedUserAccount: in.SharedUserAccount,
		SharedUserID:      in.SharedUserID,
		ProjectID:         in.ProjectID,
		NormalPerm:        in.NormalPerm,
		SystemPerm:        in.SystemPerm,
	}
}
func ToSharesTypes(in []*dm.UserDeviceShareInfo) (ret []*types.UserDeviceShareInfo) {
	for _, v := range in {
		ret = append(ret, ToShareTypes(v))
	}
	return
}
