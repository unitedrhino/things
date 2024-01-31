package dgdirect

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
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
	svcCtx       *svc.ServiceContext
	svcOnce      sync.Once
	postInitOnce sync.Once
	runSvrOnce   sync.Once
	c            config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/dg.yaml", &c)
		svcCtx = svc.NewServiceContext(c)

		startup.Init(svcCtx)
		logx.Infof("enabled dgsvr")
	})

	// 让 svcCtx 先返回，延迟执行 mqtt client 的 connect 操作.
	utils.Go(context.Background(), func() {
		postInitOnce.Do(func() {
			startup.PostInit(svcCtx)
		})
	})

	return svcCtx
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
