package config

import (
	"gitee.com/i-Things/core/shared/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Event      conf.EventConf //和things内部交互的设置
	Database   conf.Database
	CacheRedis cache.ClusterConf
	DmRpc      conf.RpcClientConf `json:",optional"`
}
