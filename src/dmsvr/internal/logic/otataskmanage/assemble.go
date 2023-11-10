package otataskmanagelogic

import (
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToOtaTaskInfo(di *relationDB.DmOtaTask) *dm.OtaTaskInfo {
	return &dm.OtaTaskInfo{
		TaskID:      di.ID,
		FirmwareID:  di.FirmwareID,
		TaskUid:     di.TaskUid,
		Type:        int32(di.Type),
		UpgradeType: int32(di.UpgradeType),
		AutoRepeat:  int32(di.AutoRepeat),
		Status:      int32(di.Status),
		DeviceList:  &wrapperspb.StringValue{Value: di.DeviceList},
		VersionList: &wrapperspb.StringValue{Value: di.VersionList},
		CreatedTime: di.CreatedTime.Unix(),
	}
}
func ToOtaTaskDeviceInfo(di *relationDB.DmOtaTaskDevices) *dm.OtaTaskDeviceInfo {
	return &dm.OtaTaskDeviceInfo{
		ID:          di.ID,
		FirmwareID:  di.FirmwareID,
		DeviceName:  di.DeviceName,
		ProductName: "", //TODO
		ProductID:   di.ProductID,
		Status:      int32(di.Status),
		Version:     di.Version,
		Desc:        di.Desc,
		RetryCount:  di.RetryCount,
		UpdatedTime: di.UpdatedTime.Unix(),
	}
}
func ToOtaTaskReadResp(di *relationDB.DmOtaTask) *dm.OtaTaskReadResp {
	return &dm.OtaTaskReadResp{
		TaskID:     di.ID,
		FirmwareID: di.FirmwareID,
		//ProductID:   di.ProductID,
		TaskUid:     di.TaskUid,
		Type:        int32(di.Type),
		UpgradeType: int32(di.UpgradeType),
		AutoRepeat:  int32(di.AutoRepeat),
		Status:      int32(di.Status),
		DeviceList:  &wrapperspb.StringValue{Value: di.DeviceList},
		VersionList: &wrapperspb.StringValue{Value: di.VersionList},
		CreatedTime: di.CreatedTime.Unix(),
	}
}
