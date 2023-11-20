package config

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Database   conf.Database
	CacheRedis cache.ClusterConf
	UserToken  struct {
		AccessSecret string
		AccessExpire int64
	}
	Register struct {
		NeedDetail   bool   `json:",default=true"` //注册的时候是否需要填写用户信息,账号密码
		DefaultRole  int64  `json:",default=1"`    //注册后的默认角色
		SecondSecret string //第二步需要的token秘钥
		SecondExpire int64  //token过期时间 单位:秒
	} `json:",optional"`
	WxMiniProgram clients.WxMiniProgram `json:",optional"` // 微信小程序，可选
	UserOpt       struct {
		NeedUserName bool  `json:",default=true"` //注册是否必须填写账号密码
		NeedPassWord bool  `json:",default=true"` //注册是否必须填写账号密码
		PassLevel    int32 `json:",default=2"`    //用户密码强度级别
	} // 用户登录注册选项
	Map struct {
		Mode      string `json:",default=baidu"`
		AccessKey string
	}
	WrongPasswordCounter conf.WrongPasswordCounter `json:",optional"`
}
