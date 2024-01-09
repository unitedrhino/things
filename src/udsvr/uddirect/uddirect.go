package uddirect

import (
	"fmt"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/udsvr/internal/config"
	ruleServer "github.com/i-Things/things/src/udsvr/internal/server/rule"
	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"
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
	ctxSvc     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/ud.yaml", &c)
		ctxSvc = svc.NewServiceContext(c)
		logx.Infof("enabled udsvr")
	})
	return ctxSvc
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
		ud.RegisterRuleServer(grpcServer, ruleServer.NewRuleServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor, ctxs.GrpcInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
