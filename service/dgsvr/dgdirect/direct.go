package dgdirect

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/interceptors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/config"
	deviceauthServer "gitee.com/unitedrhino/things/service/dgsvr/internal/server/deviceauth"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/startup"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dgsvr/pb/dg"
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
		utils.ConfMustLoad("etc/dg.yaml", &c)
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
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)

	fmt.Printf("Starting dgrpc server at %s...\n", c.ListenOn)
	s.Start()
}
