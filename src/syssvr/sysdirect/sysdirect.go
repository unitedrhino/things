package sysdirect

import (
	"github.com/i-Things/things/src/syssvr/internal/config"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"sync"
)

type Config = config.Config

var (
	ctxSvc *svc.ServiceContext
	once   sync.Once
	c      config.Config
)

func GetCtxSvc() *svc.ServiceContext {
	once.Do(func() {
		conf.MustLoad("etc/sys.yaml", &c)
		ctxSvc = svc.NewServiceContext(c)
	})
	return ctxSvc
}
