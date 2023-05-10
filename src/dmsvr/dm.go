//设备管理模块-dmsvr
package main

import (
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	deviceauth "github.com/i-Things/things/src/dmsvr/internal/server/deviceauth"
	devicemanage "github.com/i-Things/things/src/dmsvr/internal/server/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/internal/server/productmanage"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "net/http/pprof"
)

func main() {
	svcCtx := dmdirect.GetCtxSvc()
	//grpc服务初始化
	c := svcCtx.Config
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
