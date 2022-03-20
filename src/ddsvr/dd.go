package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/event"
	"github.com/i-Things/things/src/ddsvr/internal/repo/third/devLink"
	"github.com/i-Things/things/src/ddsvr/internal/repo/third/innerLink"
	"github.com/i-Things/things/src/ddsvr/internal/server"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/i-Things/things/src/ddsvr/pb/dd"
	"github.com/zeromicro/go-zero/core/logx"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dd.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	err := svcCtx.DevLink.SubScribe(func(ctx context.Context) devLink.DevSubHandle {
		return event.NewDeviceSubServer(svcCtx, ctx)
	})
	if err != nil {
		logx.Error("DevLink Subscribe failure", err)
		os.Exit(-1)
	}
	err = svcCtx.InnerLink.Subscribe(func(ctx context.Context) innerLink.InnerSubHandle {
		return event.NewInnerSubServer(svcCtx, ctx)
	})
	if err != nil {
		logx.Error("InnerLink Subscribe failure", err)
		os.Exit(-1)
	}
	svr := server.NewDdServer(svcCtx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dd.RegisterDdServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
