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
	Mongo      struct {
		Url      string //mongodb连接串
		Database string //选择的数据库
	}
	InnerLink InnerLinkConf //和things内部交互的设置
	AuthWhite AuthWhite     //设备登录校验将things服务连接的ip设备root权限
}

type InnerLinkConf struct {
	Nats conf.NatsConf
}

type AuthWhite struct {
	IpRange []string //ip 及ip段
}
