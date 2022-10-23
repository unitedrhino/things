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
	Rej        struct {
		AccessSecret string
		AccessExpire int64
	} //注册token相关配置
	FrontDir string `json:",default=./dist"` //前端文件路径
	Captcha  Captcha
	OSS      conf.OSSConf `json:",optional"`
}
