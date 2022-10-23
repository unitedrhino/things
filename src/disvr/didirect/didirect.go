package didirect

import (
	"flag"
	"github.com/i-Things/things/src/disvr/internal/config"
	"github.com/i-Things/things/src/disvr/internal/startup"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"sync"
)

type Config = config.Config

var (
	svcCtx *svc.ServiceContext
	once   sync.Once
	c      config.Config
)

func GetCtxSvc() *svc.ServiceContext {
	once.Do(func() {
		flag.Parse()
		conf.MustLoad("etc/di.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Subscribe(svcCtx)
	})
	return svcCtx
}
