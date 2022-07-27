package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/dcsvr/dcdirect"
	"github.com/i-Things/things/src/ddsvr/dddirect"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/usersvr/userdirect"
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
	UserRpc    conf.RpcClientConf `json:",optional"`
	DcRpc      conf.RpcClientConf `json:",optional"`
	DmRpc      conf.RpcClientConf `json:",optional"`
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	Rej struct {
		AccessSecret string
		AccessExpire int64
	}
	Captcha Captcha
	OSS     conf.OSSConf `json:",optional"`
}

type Configs struct {
	Config
	UserSvr userdirect.Config `json:",optional"` //只有单体模式需要填写
	DmSvr   dmdirect.Config   `json:",optional"` //只有单体模式需要填写
	DdSvr   *dddirect.Config  `json:",optional"` //只有单体模式需要填写
	DcSvr   dcdirect.Config   `json:",optional"` //只有单体模式需要填写
}
