package svc

import (
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/config"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/repo/event/publish/pubJob"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type ServiceContext struct {
	Config config.Config
	PubJob *pubJob.PubJob
}

func NewServiceContext(c config.Config) *ServiceContext {
	pj, err := pubJob.NewPubJob(c.Event)
	if err != nil {
		logx.Error("初始化消息队列 err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config: c,
		PubJob: pj,
	}
}
