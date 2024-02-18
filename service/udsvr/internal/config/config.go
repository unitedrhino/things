package config

import (
	"gitee.com/i-Things/share/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Database    conf.Database
	CacheRedis  cache.ClusterConf
	TimedJobRpc conf.RpcClientConf `json:",optional"`
	SysRpc      conf.RpcClientConf `json:",optional"`
	DmRpc       conf.RpcClientConf `json:",optional"`
	Event       conf.EventConf     //和things内部交互的设置
}
