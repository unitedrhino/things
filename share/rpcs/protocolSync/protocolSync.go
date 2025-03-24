package main

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/interceptors"

	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/config"
	protocolsyncServer "gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/server/protocolsync"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/startup"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/svc"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/pb/protocolSync"

	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	defer utils.Recover(context.Background())
	var c config.Config
	utils.ConfMustLoad("etc/protocolSync.yaml", &c)
	svcCtx := svc.NewServiceContext(c)
	startup.Init(svcCtx)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		protocolSync.RegisterProtocolSyncServer(grpcServer, protocolsyncServer.NewProtocolSyncServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
