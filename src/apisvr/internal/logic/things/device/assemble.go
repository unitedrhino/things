package device

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func DeviceInfoToApi(v *dm.DeviceInfo) *types.DeviceInfo {
	position := &types.Point{
		Longitude: v.Position.Longitude, //经度
		Latitude:  v.Position.Latitude,  //维度
	}
	return &types.DeviceInfo{
		ProductID:   v.ProductID,                   //产品id 只读
		DeviceName:  v.DeviceName,                  //设备名称 读写
		Secret:      v.Secret,                      //设备秘钥 只读
		Cert:        v.Cert,                        //设备证书 只读
		Imei:        v.Imei,                        //IMEI号信息 只读
		Mac:         v.Mac,                         //MAC号信息 只读
		Version:     utils.ToNullString(v.Version), //固件版本 读写
		HardInfo:    v.HardInfo,                    //模组硬件型号 只读
		SoftInfo:    v.SoftInfo,                    //模组软件版本 只读
		Position:    position,                      //设别定位（百度坐标）
		Address:     &v.Address.Value,              //详细地址
		Tags:        logic.ToTagsType(v.Tags),      //设备标签
		IsOnline:    v.IsOnline,                    //在线状态 1离线 2在线 只读
		FirstLogin:  v.FirstLogin,                  //激活时间 只读
		LastLogin:   v.LastLogin,                   //最后上线时间 只读
		LogLevel:    v.LogLevel,                    //日志级别 1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		CreatedTime: v.CreatedTime,                 //创建时间 只读
	}
}
