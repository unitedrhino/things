package stockdirect

import (
	"fmt"
	"github.com/i-Things/things/src/stocksvr/internal/config"
	"github.com/i-Things/things/src/stocksvr/internal/server"
	"github.com/i-Things/things/src/stocksvr/internal/startup"
	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/i-Things/things/src/stocksvr/types/pb/stock"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

type Config = config.Config

var (
	ctxSvc     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/stock.yaml", &c)
		ctxSvc = svc.NewServiceContext(c)
		startup.Init(ctxSvc)
	})
	return ctxSvc
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		go func() {
			c := svcCtx.Config
			s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
				stock.RegisterStockServer(grpcServer, server.NewStockServer(svcCtx))

				if c.Mode == service.DevMode || c.Mode == service.TestMode {
					reflection.Register(grpcServer)
				}
			})
			defer s.Stop()

			fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
			s.Start()
		}()
	})

}
