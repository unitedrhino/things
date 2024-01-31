package productCustom

import "gitee.com/i-Things/core/shared/devices"

type CustomTopic struct {
	Topic     string `json:"topic"`     //自定义主题需要以$custom 并包含设备名称{deviceName}及产品名称{productID}
	Direction int64  `json:"direction"` //1:上行 2:下行 3:双向
}

const (
	LoginAuthFuncName = "LoginAuth"
)

type LoginAuthFunc func(dir devices.Direction, clientID string, userName string, password string) (*devices.Core, error)
