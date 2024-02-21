package info

import (
	"github.com/i-Things/things/service/apisvr/internal/logic/things"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToInfoPb(in *types.DeviceTimerInfo) *ud.DeviceTimerInfo {
	if in == nil {
		return nil
	}
	return &ud.DeviceTimerInfo{
		Id:          in.ID,
		Name:        in.Name,
		Device:      things.ToUdDeviceCorePb(&in.Device),
		CreatedTime: in.CreatedTime,
		TriggerType: in.TriggerType,
		ExecAt:      in.ExecAt,
		ExecRepeat:  in.ExecRepeat,
		ActionType:  in.ActionType,
		DataID:      in.DataID,
		Value:       in.Value,
		Status:      in.Status,
		LastRunTime: in.LastRunTime,
	}
}
func ToInfoTypes(in *ud.DeviceTimerInfo) *types.DeviceTimerInfo {
	if in == nil {
		return nil
	}
	return &types.DeviceTimerInfo{
		ID:          in.Id,
		Name:        in.Name,
		Device:      *things.UdToDeviceCoreTypes(in.Device),
		CreatedTime: in.CreatedTime,
		TriggerType: in.TriggerType,
		ExecAt:      in.ExecAt,
		ExecRepeat:  in.ExecRepeat,
		ActionType:  in.ActionType,
		DataID:      in.DataID,
		Value:       in.Value,
		Status:      in.Status,
		LastRunTime: in.LastRunTime,
	}
}
func ToInfosTypes(in []*ud.DeviceTimerInfo) (ret []*types.DeviceTimerInfo) {
	for _, v := range in {
		ret = append(ret, ToInfoTypes(v))
	}
	return
}
