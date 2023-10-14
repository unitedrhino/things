package timedDelay

import (
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/clients"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"time"
)

type Option struct {
	//以下两个参数优先使用ProcessIn
	ProcessIn time.Duration //多久之后发
	ProcessAt time.Time     // 固定时间发
	Timeout   time.Duration //超时时间 优先使用
	Deadline  time.Time     //截止时间
}

type Timed struct {
	asynqClient *asynq.Client
	serverName  string //服务名
}

func NewTimed(c cache.ClusterConf, serverName string) *Timed {
	return &Timed{asynqClient: clients.NewAsynqClient(c), serverName: serverName}
}
