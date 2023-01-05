package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	TDengine struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	Event      conf.EventConf //和things内部交互的设置
	AuthWhite  conf.AuthConf  //设备登录校验将things服务连接的ip设备root权限
}
