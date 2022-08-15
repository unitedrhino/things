package main

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/src/disvr/internal/startup"

	"github.com/i-Things/things/src/disvr/internal/config"
	deviceinteractServer "github.com/i-Things/things/src/disvr/internal/server/deviceinteract"
	devicelogServer "github.com/i-Things/things/src/disvr/internal/server/devicelog"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/di.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	startup.Subscribe(svcCtx)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		di.RegisterDeviceLogServer(grpcServer, devicelogServer.NewDeviceLogServer(svcCtx))
		di.RegisterDeviceInteractServer(grpcServer, deviceinteractServer.NewDeviceInteractServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
