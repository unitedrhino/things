package ruledirect

import (
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
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
		conf.MustLoad("etc/rule.yaml", &c)
		ctxSvc = svc.NewServiceContext(c)
	})
	return ctxSvc
}
