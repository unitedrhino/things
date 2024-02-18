package ruledirect

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/interceptors"
	"github.com/i-Things/things/service/rulesvr/internal/config"
	alarmcenter "github.com/i-Things/things/service/rulesvr/internal/server/alarmcenter"
	ruleengine "github.com/i-Things/things/service/rulesvr/internal/server/ruleengine"
	scenelinkage "github.com/i-Things/things/service/rulesvr/internal/server/scenelinkage"
	"github.com/i-Things/things/service/rulesvr/internal/startup"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/internal/timer/sceneTimer"
	"github.com/i-Things/things/service/rulesvr/pb/rule"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
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
		startup.Init(svcCtx)
		sceneTimer.NewSceneTimer(context.TODO(), svcCtx).Start()
		svcCtx.SceneTimerControl = sceneTimer.NewSceneTimerControl()
		logx.Infof("enabled rulesvr")
	})
	return svcCtx
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		go Run(svcCtx)
	})

}

func Run(svcCtx *svc.ServiceContext) {
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
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)

	fmt.Printf("Starting rulesvr server at %s...\n", c.ListenOn)
	s.Start()
}
