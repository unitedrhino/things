package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Event     conf.EventConf //和things内部交互的设置
	SqlConf   SqlConf
	CacheConf cache.CacheConf
	DiRpc     conf.RpcClientConf `json:",optional"`
}
type SqlConf struct {
	DataSource string
}
