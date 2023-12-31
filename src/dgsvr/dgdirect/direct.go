package dgdirect

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dgsvr/internal/config"
	deviceauthServer "github.com/i-Things/things/src/dgsvr/internal/server/deviceauth"
	"github.com/i-Things/things/src/dgsvr/internal/startup"
	"github.com/i-Things/things/src/dgsvr/internal/svc"
	"github.com/i-Things/things/src/dgsvr/pb/dg"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

type Config = config.Config

var (
	svcCtx     svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/dg.yaml", &c)
		utils.Go(context.Background(), func() {
			svcCtxNew := svc.NewServiceContext(c)
			startup.Init(svcCtxNew)
			svcCtx = *svcCtxNew
			logx.Infof("enabled dgsvr")
		})
	})
	return &svcCtx
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		go Run(svcCtx)
	})

}
func Run(svcCtx *svc.ServiceContext) {
	c := svcCtx.Config
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dg.RegisterDeviceAuthServer(grpcServer, deviceauthServer.NewDeviceAuthServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor, ctxs.GrpcInterceptor)

	fmt.Printf("Starting dgrpc server at %s...\n", c.ListenOn)
	s.Start()
}
