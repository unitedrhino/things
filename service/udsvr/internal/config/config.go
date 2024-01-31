package config

import (
	"gitee.com/i-Things/share/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Database   conf.Database
	CacheRedis cache.ClusterConf
}
