package config

import (
	"gitee.com/i-Things/core/shared/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DmRpc   conf.RpcClientConf `json:",optional"`
	DevLink conf.DevLinkConf   //和设备交互的设置
	Event   conf.EventConf     //和things内部交互的设置
}
