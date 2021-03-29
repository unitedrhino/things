package config

import (
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	UserToken struct {
		AccessSecret string
		AccessExpire int64
	}
	Rej struct{
		AccessSecret string
		AccessExpire int64
	}
	NodeID int64
}
