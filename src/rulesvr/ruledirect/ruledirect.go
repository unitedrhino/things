package ruledirect

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/config"
	alarmcenter "github.com/i-Things/things/src/rulesvr/internal/server/alarmcenter"
	ruleengine "github.com/i-Things/things/src/rulesvr/internal/server/ruleengine"
	scenelinkage "github.com/i-Things/things/src/rulesvr/internal/server/scenelinkage"
	"github.com/i-Things/things/src/rulesvr/internal/startup"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/internal/timer/sceneTimer"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"sync"
)

type Config = config.Config

var (
	svcCtx     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
	ConfigFile = "etc/rule.yaml"
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad(ConfigFile, &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Subscribe(svcCtx)
		sceneTimer.NewSceneTimer(context.TODO(), svcCtx).Start()
		svcCtx.SceneTimerControl = sceneTimer.NewSceneTimerControl()
	})
	return svcCtx
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		go func() {
			c := svcCtx.Config

			s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
				rule.RegisterSceneLinkageServer(grpcServer, scenelinkage.NewSceneLinkageServer(svcCtx))
				rule.RegisterRuleEngineServer(grpcServer, ruleengine.NewRuleEngineServer(svcCtx))
				rule.RegisterAlarmCenterServer(grpcServer, alarmcenter.NewAlarmCenterServer(svcCtx))
				if c.Mode == service.DevMode || c.Mode == service.TestMode {
					reflection.Register(grpcServer)
				}
			})
			defer s.Stop()
			s.AddUnaryInterceptors(errors.ErrorInterceptor)

			fmt.Printf("Starting rulesvr server at %s...\n", c.ListenOn)
			s.Start()
		}()
	})

}
