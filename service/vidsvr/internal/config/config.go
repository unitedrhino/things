package config

import (
	"gitee.com/i-Things/share/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Database   conf.Database
	CacheRedis cache.ClusterConf
	Event      conf.EventConf //和things内部交互的设置
	Mediakit   conf.MediaConf `json:",optional"` //docker -zlemdiakit的连接
	Restconf   rest.RestConf  //docker访问
}
