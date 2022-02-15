package weixin

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-zero/core/logx"
	zeroCache "github.com/zeromicro/go-zero/core/stores/cache"
)

type MiniprogramConf struct {
	Open      bool //如果开启则需要初始化为true
	AppID     string
	AppSecret string
}

type MiniProgram = miniprogram.MiniProgram

func NewWexinMiniProgram(conf MiniprogramConf, redisConf zeroCache.ClusterConf) *MiniProgram {
	if conf.Open == false {
		logx.Info("weixin mini program conf is not open")
		return nil
	}
	wc := wechat.NewWechat()
	memory := cache.NewRedis(&cache.RedisOpts{
		Host:     redisConf[0].Host,
		Password: redisConf[0].Pass,
	})
	cfg := &miniConfig.Config{
		AppID:     conf.AppID,
		AppSecret: conf.AppSecret,
		Cache:     memory,
	}
	miniprogram := wc.GetMiniProgram(cfg)
	return miniprogram
}
