package main

import (
	"flag"
	"fmt"

	"github.com/i-Things/things/src/vidsvr/internal/config"
	vidmgrmangeServer "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrmange"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/vid.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		vid.RegisterVidmgrMangeServer(grpcServer, vidmgrmangeServer.NewVidmgrMangeServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			fmt.Printf("Starting rpc serverc.Mode == service.DevMode\n")
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
