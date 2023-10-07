package main

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/src/timedqueuesvr/scheduler/internal/startup"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/i-Things/things/src/timedqueuesvr/scheduler/internal/config"
	schedulerServer "github.com/i-Things/things/src/timedqueuesvr/scheduler/internal/server/scheduler"
	"github.com/i-Things/things/src/timedqueuesvr/scheduler/internal/svc"
	"github.com/i-Things/things/src/timedqueuesvr/scheduler/pb/scheduler"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/scheduler.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	startup.Init(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		scheduler.RegisterSchedulerServer(grpcServer, schedulerServer.NewSchedulerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
