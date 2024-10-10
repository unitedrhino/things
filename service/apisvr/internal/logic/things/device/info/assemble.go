package info

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func ToRpcDeviceInfo(req *types.DeviceInfo) *dm.DeviceInfo {
	return utils.Copy[dm.DeviceInfo](req)
}
