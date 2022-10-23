package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func toTagsMap(tags []*types.DeviceTag) map[string]string {
	if tags == nil {
		return nil
	}
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return tagMap
}
func toTagsType(tags map[string]string) (retTag []*types.DeviceTag) {
	for k, v := range tags {
		retTag = append(retTag, &types.DeviceTag{
			Key:   k,
			Value: v,
		})
	}
	return
}

func deviceInfoToApi(v *dm.DeviceInfo) *types.DeviceInfo {
	return &types.DeviceInfo{
		ProductID:   v.ProductID,                   //产品id 只读
		DeviceName:  v.DeviceName,                  //设备名称 读写
		CreatedTime: v.CreatedTime,                 //创建时间 只读
		Secret:      v.Secret,                      //设备秘钥 只读
		FirstLogin:  v.FirstLogin,                  //激活时间 只读
		LastLogin:   v.LastLogin,                   //最后上线时间 只读
		Version:     utils.ToNullString(v.Version), // 固件版本  读写
		LogLevel:    v.LogLevel,                    // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		Cert:        v.Cert,                        // 设备证书  只读
		Tags:        toTagsType(v.Tags),
		IsOnline:    v.IsOnline, // 在线状态  1离线 2在线 只读
	}
}
