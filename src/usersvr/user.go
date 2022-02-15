package main

import (
	"flag"
	"fmt"
	"github.com/go-things/things/shared/errors"
	"google.golang.org/grpc/reflection"

	"github.com/go-things/things/src/usersvr/internal/config"
	"github.com/go-things/things/src/usersvr/internal/server"
	"github.com/go-things/things/src/usersvr/internal/svc"
	"github.com/go-things/things/src/usersvr/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewUserServer(ctx)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, srv)
		reflection.Register(grpcServer)
	})
	s.AddUnaryInterceptors(errors.ErrorInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
