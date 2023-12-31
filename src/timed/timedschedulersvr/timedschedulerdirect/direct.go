package timedschedulerdirect

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/config"
	schedulerServer "github.com/i-Things/things/src/timed/timedschedulersvr/internal/server/timedscheduler"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/startup"
	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedschedulersvr/pb/timedscheduler"
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
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		flag.Parse()
		conf.MustLoad("etc/timedscheduler.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Init(svcCtx)
		logx.Infof("enabled timedschedulersvr")
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
		timedscheduler.RegisterTimedschedulerServer(grpcServer, schedulerServer.NewTimedschedulerServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor, ctxs.GrpcInterceptor)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
