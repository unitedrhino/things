package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/event/eventDevSub"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/innerLink"
	"github.com/i-Things/things/src/dmsvr/internal/server"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/internal/vars"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
	svcCtx := svc.NewServiceContext(c)
	vars.Svrctx = svcCtx

	svcCtx.InnerLink.Subscribe(func(ctx context.Context) innerLink.InnerSubHandle {
		return eventDevSub.NewDeviceMsgHandle(ctx, svcCtx)
	})
	//grpc服务初始化

	srv := server.NewDmServer(svcCtx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dm.RegisterDmServer(grpcServer, srv)
		reflection.Register(grpcServer)
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
