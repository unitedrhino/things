package userdevicelogic

import (
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

const deviceShareUseByWechatSingleDevice = "wechat_single_device"

func shouldConsumeShareTokenAfterAccept(useBy string) bool {
	return useBy == deviceShareUseByWechatSingleDevice
}

func buildMultiShareTokenResponse(shareToken string, info *dm.UserDeviceShareMultiInfo) *dm.UserDeviceShareMultiToken {
	resp := &dm.UserDeviceShareMultiToken{ShareToken: shareToken}
	if info == nil {
		return resp
	}
	resp.CreatedTime = info.CreatedTime
	resp.LinkExpireAt = calcMultiShareLinkExpireAt(info.CreatedTime)
	resp.AuthExpireAt = info.ExpTime
	return resp
}

func calcMultiShareLinkExpireAt(createdTime int64) int64 {
	if createdTime <= 0 {
		return 0
	}
	return createdTime + userShared.MultiDeviceShareTokenTTLSeconds
}
