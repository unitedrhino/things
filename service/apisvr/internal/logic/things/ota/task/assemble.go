package task

import (
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func otaTaskDeviceInfoToApi(v *dm.OtaTaskDeviceInfo) *types.OtaTaskDeviceInfo {
	return &types.OtaTaskDeviceInfo{
		ID:          v.ID,
		TaskUid:     v.TaskUid,
		Version:     v.Version,
		DeviceName:  v.DeviceName,
		Status:      v.Status,
		UpdatedTime: v.UpdatedTime,
	}
}
func otaTaskInfoToApi(v *dm.OtaTaskInfo) *types.OtaTaskInfo {
	return &types.OtaTaskInfo{
		TaskID:      v.TaskID,
		TaskUid:     v.TaskUid,
		Type:        v.Type,
		UpgradeType: v.UpgradeType,
		Status:      v.Status,
		CreatedTime: v.CreatedTime,
	}
}
func otaTaskReadToApi(v *dm.OtaTaskReadResp) *types.OtaTaskReadResp {
	return &types.OtaTaskReadResp{
		TaskID:      v.TaskID,
		TaskUid:     v.TaskUid,
		Type:        v.Type,
		UpgradeType: v.UpgradeType,
		Version:     "", //TODO
		SrcVersion:  v.VersionList.GetValue(),
		SrcDevice:   v.DeviceList.GetValue(),
		ProductID:   "", //TODO
		ProductName: "", //TODO
		AutoRepeat:  v.AutoRepeat,
		Status:      v.Status,
		CreatedTime: v.CreatedTime,
	}
}
