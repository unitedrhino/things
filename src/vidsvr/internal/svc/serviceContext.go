package svc

import (
	"fmt"
	"github.com/i-Things/things/shared/eventBus"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	"github.com/i-Things/things/src/vidsvr/internal/media"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config   config.Config
	VidmgrID *utils.SnowFlake
	Cache    kv.Store
	Bus      eventBus.Bus
}

func NewServiceContext(c config.Config) *ServiceContext {

	//pd, err := pubDev
	cache := kv.NewStore(c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	VidmgrID := utils.NewSnowFlake(nodeId)
	bus := eventBus.NewEventBus()
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("vidsvr 数据库初始化失败 err", err)
		os.Exit(-1)
	} else {
		fmt.Printf("Vidsvr 数据库初始化成功 \n")
	}

	media.NewMediaChan(c)

	svcCtx := &ServiceContext{
		Config:   c,
		VidmgrID: VidmgrID,
		Cache:    cache,
		Bus:      bus,
	}
	return svcCtx
}
