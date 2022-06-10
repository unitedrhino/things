package main

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/src/filesvr/internal/handler"
	"github.com/zeromicro/go-zero/rest"
	"sync"

	"github.com/i-Things/things/src/filesvr/internal/config"
	"github.com/i-Things/things/src/filesvr/internal/server"
	"github.com/i-Things/things/src/filesvr/internal/svc"
	"github.com/i-Things/things/src/filesvr/pb/file"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/file.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	var wait sync.WaitGroup
	wait.Add(1)
	go RunRpc(c, ctx, &wait)
	wait.Add(1)
	go RunHttp(c, ctx, &wait)
	wait.Wait()
}

func RunHttp(c config.Config, ctx *svc.ServiceContext, wait *sync.WaitGroup) {
	defer wait.Done()
	server := rest.MustNewServer(c.Rest)
	defer server.Stop()
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Rest.Host, c.Rest.Port)
	server.Start()
}

func RunRpc(c config.Config, ctx *svc.ServiceContext, wait *sync.WaitGroup) {
	defer wait.Done()
	svr := server.NewFileServer(ctx)

	s := zrpc.MustNewServer(c.Rpc, func(grpcServer *grpc.Server) {
		file.RegisterFileServer(grpcServer, svr)
		if c.Rpc.Mode == service.DevMode || c.Rpc.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.Rpc.ListenOn)
	s.Start()
}
