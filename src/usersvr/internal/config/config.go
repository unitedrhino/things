package config

import (
	"github.com/i-Things/things/shared/third/weixin"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	UserToken  struct {
		AccessSecret string
		AccessExpire int64
	}
	Rej struct {
		AccessSecret string
		AccessExpire int64
	}
	WexinMiniprogram weixin.MiniprogramConf `json:",optional"` // 微信小程序，可选
	UserOpt          struct {
		NeedUserName bool  `json:default=true` //注册是否必须填写账号密码
		NeedPassWord bool  `json:default=true` //注册是否必须填写账号密码
		PassLevel    int32 `json:default=2`    //用户密码强度级别
	} // 用户登录注册选项
}
