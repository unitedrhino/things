package didirect

import (
	"github.com/i-Things/things/src/disvr/internal/config"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"sync"
)

type Config = config.Config

var (
	ctxSvc *svc.ServiceContext
	once   sync.Once
)

func getCtxSvc(config *Config) *svc.ServiceContext {
	once.Do(func() {
		ctxSvc = svc.NewServiceContext(*config)
	})
	return ctxSvc
}
