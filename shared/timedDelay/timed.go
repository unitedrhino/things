package timedDelay

import (
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/clients"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Timed struct {
	asynqClient *asynq.Client
	serverName  string //服务名
}

func NewTimed(c cache.ClusterConf, serverName string) *Timed {
	return &Timed{asynqClient: clients.NewAsynqClient(c), serverName: serverName}
}
