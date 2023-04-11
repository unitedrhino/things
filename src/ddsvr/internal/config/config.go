package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DmRpc   conf.RpcClientConf `json:",optional"`
	DevLink conf.DevLinkConf   //和设备交互的设置
	Event   conf.EventConf     //和things内部交互的设置
}
