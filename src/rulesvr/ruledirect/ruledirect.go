package ruledirect

import (
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/startup"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"sync"
)

type Config = config.Config

var (
	svcCtx *svc.ServiceContext
	once   sync.Once
	c      config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	once.Do(func() {
		conf.MustLoad("etc/rule.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Subscribe(svcCtx)
	})
	return svcCtx
}
