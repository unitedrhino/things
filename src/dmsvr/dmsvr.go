package main

import (
	"flag"
	"fmt"
	"yl/src/dmsvr/device"
	"yl/src/dmsvr/dm"
	"yl/src/dmsvr/internal/config"
	"yl/src/dmsvr/internal/server"
	"yl/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/dmsvr.yaml", "the config file")

func main() {
	device.Start()

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewDmServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dm.RegisterDmServer(grpcServer, srv)
	})
	defer s.Stop()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
