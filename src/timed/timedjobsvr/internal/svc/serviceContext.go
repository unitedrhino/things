package svc

import (
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timed/internal/repo/event/publish/pubJob"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedjobsvr/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
)

type ServiceContext struct {
	Config      config.Config
	Store       kv.Store
	Redis       *redis.Redis
	PubJob      *pubJob.PubJob
	AsynqClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	pj, err := pubJob.NewPubJob(c.Event)
	if err != nil {
		logx.Error("初始化消息队列 err", err)
		os.Exit(-1)
	}
	stores.InitConn(c.Database)
	err = relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("timedjobsvr 数据库初始化失败 err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:      c,
		PubJob:      pj,
		Redis:       redis.MustNewRedis(c.CacheRedis[0].RedisConf),
		AsynqClient: clients.NewAsynqClient(c.CacheRedis),
		Store:       kv.NewStore(c.CacheRedis),
	}
}
