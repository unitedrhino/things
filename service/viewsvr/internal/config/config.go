package config

import (
	"gitee.com/i-Things/core/shared/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database   conf.Database
	CacheRedis cache.ClusterConf
	OssConf    conf.OssConf `json:",optional"`
}
