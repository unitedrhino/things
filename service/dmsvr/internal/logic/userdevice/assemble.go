package userdevicelogic

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToUserDeviceSharePb(in *relationDB.DmUserDeviceShare) *dm.UserDeviceShareInfo {
	if in == nil {
		return nil
	}
	return &dm.UserDeviceShareInfo{
		Id:        in.ID,
		ProjectID: int64(in.ProjectID),
		Device: &dm.DeviceCore{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		},
		CreatedTime:       in.CreatedTime.Unix(),
		AuthType:          in.AuthType,
		SharedUserAccount: in.SharedUserAccount,
		SharedUserID:      in.SharedUserID,
		AccessPerm:        utils.CopyMap[dm.SharePerm](in.AccessPerm),
		SchemaPerm:        utils.CopyMap[dm.SharePerm](in.SchemaPerm),
	}
}
func ToUserDeviceSharePbs(in []*relationDB.DmUserDeviceShare) (ret []*dm.UserDeviceShareInfo) {
	for _, v := range in {
		ret = append(ret, ToUserDeviceSharePb(v))
	}
	return
}
