package main

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/config"
	jobServer "github.com/i-Things/things/src/timedqueuesvr/job/internal/server/job"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/startup"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/svc"
	"github.com/i-Things/things/src/timedqueuesvr/job/pb/job"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/job.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	startup.Init(ctx)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		job.RegisterJobServer(grpcServer, jobServer.NewJobServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
