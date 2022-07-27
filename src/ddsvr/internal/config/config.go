package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DevLink   DevLinkConf        //和设备交互的设置
	InnerLink conf.InnerLinkConf //和things内部交互的设置
}
type DevLinkConf struct {
	Mode    string         `json:",default=mqtt"` //模式 默认mqtt
	SubMode string         `json:",default=emq"`  //
	Mqtt    *conf.MqttConf `json:",optional"`
}
