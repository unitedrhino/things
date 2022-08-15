package main

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	deviceauth "github.com/i-Things/things/src/dmsvr/internal/server/deviceauth"
	devicemanage "github.com/i-Things/things/src/dmsvr/internal/server/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/internal/server/productmanage"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "net/http/pprof"
)

var configFile = flag.String("f", "etc/dm.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	//grpc服务初始化
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dm.RegisterDeviceAuthServer(grpcServer, deviceauth.NewDeviceAuthServer(svcCtx))
		dm.RegisterDeviceManageServer(grpcServer, devicemanage.NewDeviceManageServer(svcCtx))
		dm.RegisterProductManageServer(grpcServer, productmanage.NewProductManageServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
