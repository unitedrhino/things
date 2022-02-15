package main

import (
	"flag"
	"fmt"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/config"
	"github.com/go-things/things/src/dmsvr/internal/exchange"
	"github.com/go-things/things/src/dmsvr/internal/server"
	"github.com/go-things/things/src/dmsvr/internal/svc"
	"github.com/go-things/things/src/dmsvr/internal/vars"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	_ "net/http/pprof"
)

var configFile = flag.String("f", "etc/dm.yaml", "the config file")

func main() {
	flag.Parse()
	//go device.NewDevice()
	//device.TestMongo()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	//kafka服务初始化
	ctx := svc.NewServiceContext(c)
	vars.Svrctx = ctx
	k := exchange.NewKafka(ctx)
	k.AddRouters()
	stop := k.Start()
	defer stop()

	//grpc服务初始化

	srv := server.NewDmServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dm.RegisterDmServer(grpcServer, srv)
		reflection.Register(grpcServer)

	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
