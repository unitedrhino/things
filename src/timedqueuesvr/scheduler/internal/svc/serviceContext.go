package svc

import (
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timedqueuesvr/scheduler/internal/config"
	"github.com/i-Things/things/src/timedqueuesvr/scheduler/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type ServiceContext struct {
	Config    config.Config
	Scheduler *asynq.Scheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	err := relationDB.Migrate()
	if err != nil {
		logx.Error("初始化数据库错误 err", err)
		os.Exit(-1)
	}
	Scheduler := newScheduler(c)
	return &ServiceContext{
		Scheduler: Scheduler,
		Config:    c,
	}
}
