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
	InnerLink  conf.InnerLinkConf //和things内部交互的设置
	AuthWhite  AuthWhite          //设备登录校验将things服务连接的ip设备root权限
}

type AuthWhite struct {
	IpRange []string //ip 及ip段
}
