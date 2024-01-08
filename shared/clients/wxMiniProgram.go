package clients

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	zeroCache "github.com/zeromicro/go-zero/core/stores/cache"
)

type MiniProgram = miniprogram.MiniProgram

func NewWxMiniProgram(ctx context.Context, conf *conf.ThirdConf, redisConf zeroCache.ClusterConf) (*MiniProgram, error) {
	if conf == nil {
		return nil, nil
	}
	wc := wechat.NewWechat()
	memory := cache.NewRedis(ctx, &cache.RedisOpts{
		Host:     redisConf[0].Host,
		Password: redisConf[0].Pass,
	})
	cfg := &miniConfig.Config{
		AppID:     conf.AppID,
		AppSecret: conf.AppSecret,
		Cache:     memory,
	}
	program := wc.GetMiniProgram(cfg)
	return program, nil
}
