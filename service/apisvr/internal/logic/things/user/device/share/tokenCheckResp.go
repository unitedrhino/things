package share

import (
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

const multiDeviceShareTokenTTLSeconds int64 = 24 * 60 * 60
const tokenCheckReasonExpiredOrConsumed = "expired_or_consumed"

func ToTokenCheckResp(in *dm.UserDeviceShareMultiInfo) *types.UserDeviceShareTokenCheckResp {
	if in == nil {
		return InvalidTokenCheckResp()
	}
	return &types.UserDeviceShareTokenCheckResp{
		Valid:        true,
		LinkExpireAt: calcTokenCheckLinkExpireAt(in.CreatedTime),
		AuthExpireAt: in.ExpTime,
		CreatedTime:  in.CreatedTime,
		UseBy:        in.UseBy,
		Desc:         in.Desc,
		DeviceCount:  int64(len(in.Devices)),
	}
}

func InvalidTokenCheckResp() *types.UserDeviceShareTokenCheckResp {
	return &types.UserDeviceShareTokenCheckResp{
		Valid:  false,
		Reason: tokenCheckReasonExpiredOrConsumed,
	}
}

func calcTokenCheckLinkExpireAt(createdTime int64) int64 {
	if createdTime <= 0 {
		return 0
	}
	return createdTime + multiDeviceShareTokenTTLSeconds
}
