package info

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToRpcDeviceInfo(req *types.DeviceInfo) *dm.DeviceInfo {
	return utils.Copy[dm.DeviceInfo](req)
}
