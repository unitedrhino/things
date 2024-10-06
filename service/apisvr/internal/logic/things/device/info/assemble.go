package info

import (
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToRpcDeviceInfo(req *types.DeviceInfo) *dm.DeviceInfo {
	return utils.Copy[dm.DeviceInfo](req)
}
