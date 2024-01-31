package info

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToRpcDeviceInfo(req *types.DeviceInfoSaveReq) *dm.DeviceInfo {
	dmReq := &dm.DeviceInfo{
		ProductID:      req.ProductID,  //产品id 只读
		DeviceName:     req.DeviceName, //设备名称 读写
		LogLevel:       req.LogLevel,   // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		Tags:           logic.ToTagsMap(req.Tags),
		Address:        utils.ToRpcNullString(req.Address),
		Position:       logic.ToDmPointRpc(req.Position),
		DeviceAlias:    utils.ToRpcNullString(req.DeviceAlias), //设备别名 读写
		Phone:          utils.ToRpcNullString(req.Phone),
		Iccid:          utils.ToRpcNullString(req.Iccid),
		MobileOperator: req.MobileOperator,
		AreaID:         req.AreaID, //项目区域id 只读
	}
	return dmReq
}
