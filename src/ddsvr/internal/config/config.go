package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DevLink   DevLinkConf   //和设备交互的设置
	InnerLink InnerLinkConf //和things内部交互的设置
}

type DevLinkConf struct {
	Mode    string         `json:",default=mqtt"` //模式 默认mqtt
	SubMode string         `json:",default=emq"`  //
	Mqtt    *conf.MqttConf `json:",optional"`
}

type InnerLinkConf struct {
	Nats conf.NatsConf
}
