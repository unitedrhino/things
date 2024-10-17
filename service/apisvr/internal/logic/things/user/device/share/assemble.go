package share

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
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
		AuthType:     in.AuthType,
		ProjectID:    in.ProjectID,
		ExpTime:      utils.ToRpcNullInt64(in.ExpTime),
		AccessPerm:   utils.CopyMap[dm.SharePerm](in.AccessPerm),
		SchemaPerm:   utils.CopyMap[dm.SharePerm](in.SchemaPerm),
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
		AuthType:          in.AuthType,
		CreatedTime:       in.CreatedTime,
		ExpTime:           utils.ToNullInt64(in.ExpTime),
		SharedUserAccount: in.SharedUserAccount,
		SharedUserID:      in.SharedUserID,
		ProjectID:         in.ProjectID,
		AccessPerm:        utils.CopyMap[types.SharePerm](in.AccessPerm),
		SchemaPerm:        utils.CopyMap[types.SharePerm](in.SchemaPerm),
	}
}
func ToSharesTypes(in []*dm.UserDeviceShareInfo) (ret []*types.UserDeviceShareInfo) {
	for _, v := range in {
		ret = append(ret, ToShareTypes(v))
	}
	return
}
func ToMuitlSharePb(in *types.UserMultiDevicesShareInfo) *dm.UserMultiDevicesShareInfo {
	if in == nil {
		return nil
	}
	devices := in.Devices
	return &dm.UserMultiDevicesShareInfo{
		Device:   toSharesDevices(devices),
		AuthType: in.AuthType,
		//ProjectID:    in.ProjectID,
		ExpTime:    in.ExpTime,
		AccessPerm: utils.CopyMap[dm.SharePerm](in.AccessPerm),
		SchemaPerm: utils.CopyMap[dm.SharePerm](in.SchemaPerm),
	}
}
func toSharesDevices(in []*types.DeviceCore) (ret []*dm.DeviceCore) {
	for _, v := range in {
		ret = append(ret, &dm.DeviceCore{
			DeviceName: v.DeviceName,
			ProductID:  v.ProductID,
		})
	}
	return ret
}
func ToMultiShareTypes(in *dm.UserMultiDevicesShareInfo) *types.UserMultiDevicesShareInfo {
	if in == nil {
		return nil
	}
	var dvs []*types.DeviceCore
	for _, v := range in.Device {
		dvs = append(dvs, &types.DeviceCore{
			DeviceName: v.DeviceName,
			ProductID:  v.ProductID,
		})
	}
	return &types.UserMultiDevicesShareInfo{
		Devices:     dvs,
		AuthType:    in.AuthType,
		CreatedTime: in.CreatedTime,
		ExpTime:     in.ExpTime,
		// SharedUserAccount: in.SharedUserAccount,
		// SharedUserID:      in.SharedUserID,
		//ProjectID:         in.ProjectID,
		AccessPerm: utils.CopyMap[types.SharePerm](in.AccessPerm),
		SchemaPerm: utils.CopyMap[types.SharePerm](in.SchemaPerm),
	}
}
