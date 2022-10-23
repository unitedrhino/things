package main

import (
	"fmt"
	"github.com/i-Things/things/src/disvr/didirect"
	deviceinteract "github.com/i-Things/things/src/disvr/internal/server/deviceinteract"
	devicemsg "github.com/i-Things/things/src/disvr/internal/server/devicemsg"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	svcCtx := didirect.GetCtxSvc()
	c := svcCtx.Config
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		di.RegisterDeviceMsgServer(grpcServer, devicemsg.NewDeviceMsgServer(svcCtx))
		di.RegisterDeviceInteractServer(grpcServer, deviceinteract.NewDeviceInteractServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
