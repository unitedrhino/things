package svc

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config       config.Config
	Scheduler    *clients.TimedScheduler
	Store        kv.Store
	SchedulerRun bool //只启动单例
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("初始化数据库错误 err", err)
		os.Exit(-1)
	}
	Scheduler := clients.NewTimedScheduler(c.CacheRedis)
	return &ServiceContext{
		Scheduler: Scheduler,
		Config:    c,
		Store:     kv.NewStore(c.CacheRedis),
	}
}
