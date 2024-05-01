package share

import (
	"gitee.com/i-Things/share/utils"
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
		CreatedTime:       in.CreatedTime,
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
