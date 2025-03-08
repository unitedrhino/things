package config

import (
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Etcd        discov.EtcdConf `json:",optional,inherit"`
	CacheRedis  cache.ClusterConf
	SysRpc      conf.RpcClientConf `json:",optional"`
	DgRpc       conf.RpcClientConf `json:",optional"`
	DmRpc       conf.RpcClientConf `json:",optional"`
	UdRpc       conf.RpcClientConf `json:",optional"`
	TimedJobRpc conf.RpcClientConf `json:",optional"`
	OssConf     conf.OssConf       `json:",optional"`
	Event       conf.EventConf     //和things内部交互的设置
}
