package svc

import (
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/eventBus"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
	"github.com/i-Things/things/src/timed/timedjobsvr/timedjobdirect"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type ServiceContext struct {
	Config   config.Config
	VidmgrID *utils.SnowFlake
	Cache    kv.Store
	Bus      eventBus.Bus
	TimedM   timedmanage.TimedManage
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedM timedmanage.TimedManage
	)

	//pd, err := pubDev
	cache := kv.NewStore(c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	VidmgrID := utils.NewSnowFlake(nodeId)
	bus := eventBus.NewEventBus()
	stores.InitConn(c.Database)
	err := relationDB.Migrate()
	if err != nil {
		logx.Error("vidsvr 数据库初始化失败 err", err)
		os.Exit(-1)
	} else {
		fmt.Printf("Vidsvr 数据库初始化成功 \n")
	}

	if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
		timedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
	} else {
		timedM = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
	}

	return &ServiceContext{
		Config:   c,
		VidmgrID: VidmgrID,
		Cache:    cache,
		Bus:      bus,
		TimedM:   timedM,
	}
}
