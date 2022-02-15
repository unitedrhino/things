package config

import (
	"github.com/go-things/things/shared/conf"
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
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	UserRpc    conf.RpcClientConf
	DcRpc      conf.RpcClientConf
	DmRpc      conf.RpcClientConf
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	Rej struct {
		AccessSecret string
		AccessExpire int64
	}
	Captcha
	NodeID int64
}
