package deviceAuth

import (
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
)

type (
	AccessAuthInfo struct {
		Username string            //用户名
		Topic    string            //主题
		ClientID string            //clientID
		Access   devices.Direction //操作
		Ip       string            //访问的ip地址
	}
)

/*
	系统topic及物模型topic都是
		第一个表示大的功能(如$thing,$ota)
		第二个表示上行还是下行
		中间为自定义字段
		以产品id/设备名结尾
*/
func AccessAuth(in AccessAuthInfo) error {
	lg, err := GetClientIDInfo(in.ClientID)
	if err != nil {
		return err
	}
	topicInfo, err := devices.GetTopicInfo(in.Topic)
	if err != nil {
		return errors.Permissions
	}
	if in.Access != topicInfo.Direction {
		return errors.Permissions
	}
	if topicInfo.ProductID != lg.ProductID || topicInfo.DeviceName != lg.DeviceName {
		return errors.Permissions
	}
	return nil
}
