package config

import (
	"gitee.com/godLei6/things/shared/conf"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	NodeID int64 //节点id
	Mysql  struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	DmRpc      conf.RpcClientConf
}
