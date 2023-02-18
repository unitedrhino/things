package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Captcha struct {
	KeyLong   int   `json:",default=6"`
	ImgWidth  int   `json:",default=240"`
	ImgHeight int   `json:",default=80"`
	KeepTime  int64 `json:",default=180"`
}

type Config struct {
	rest.RestConf
	CacheRedis cache.ClusterConf
	DdEnable   bool               `json:",optional"`
	SysRpc     conf.RpcClientConf `json:",optional"`
	DiRpc      conf.RpcClientConf `json:",optional"`
	DmRpc      conf.RpcClientConf `json:",optional"`
	RuleRpc    conf.RpcClientConf `json:",optional"`
	Rej        struct {
		AccessSecret string
		AccessExpire int64
	} //注册token相关配置
	Proxy    conf.ProxyConf
	Captcha  Captcha
	OSS      conf.OSSConf  `json:",optional"`
	Map      conf.MapConf  `json:",optional"`
	OpenAuth conf.AuthConf `json:",optional"`
}
