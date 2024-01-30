package svc

import (
	"fmt"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/vidsip/internal/config"
	"github.com/i-Things/things/src/vidsip/internal/media"
	"github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config config.Config
	Cache  kv.Store
}

func NewServiceContext(c config.Config) *ServiceContext {

	cache := kv.NewStore(c.CacheRedis)
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("vidsip 数据库初始化失败 err", err)
		os.Exit(-1)
	} else {
		fmt.Printf("Vidsip 数据库初始化成功 \n")
	}

	media.NewSipServer(c)

	return &ServiceContext{
		Config: c,
		Cache:  cache,
	}
}
