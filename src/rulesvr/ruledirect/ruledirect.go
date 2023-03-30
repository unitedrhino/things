package ruledirect

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/startup"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/internal/timer/sceneTimer"
	"github.com/zeromicro/go-zero/core/conf"

	"sync"
)

type Config = config.Config

var (
	svcCtx     *svc.ServiceContext
	once       sync.Once
	c          config.Config
	ConfigFile = "etc/rule.yaml"
)

func GetSvcCtx() *svc.ServiceContext {
	once.Do(func() {
		conf.MustLoad(ConfigFile, &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Subscribe(svcCtx)
		sceneTimer.NewSceneTimer(context.TODO(), svcCtx).Start()
		svcCtx.SceneTimerControl = sceneTimer.NewSceneTimerControl()
	})
	return svcCtx
}
