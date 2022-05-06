package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/dmsvr/internal/event/deviceMsgEvent"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/innerLink"
	"github.com/i-Things/things/src/dmsvr/internal/server"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	_ "net/http/pprof"
)

var configFile = flag.String("f", "etc/dm.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	Subscribe(svcCtx)
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

func Subscribe(svcCtx *svc.ServiceContext) {
	err := svcCtx.InnerLink.Subscribe(func(ctx context.Context) innerLink.InnerSubEvent {
		return deviceMsgEvent.NewDeviceMsgHandle(ctx, svcCtx)
	})
	if err != nil {
		log.Fatalf("%v|InnerLink.Subscribe|err:%v",
			utils.FuncName(), err)
	}
	err = svcCtx.DataUpdate.Subscribe(func(ctx context.Context) dataUpdate.DataUpdateSubEvent {
		return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	})
	if err != nil {
		log.Fatalf("%v|DataUpdate.Subscribe|err:%v",
			utils.FuncName(), err)
	}
}
