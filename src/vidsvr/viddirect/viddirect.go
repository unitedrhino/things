package viddirect

import (
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	vidmanage "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrmange"
	"github.com/i-Things/things/src/vidsvr/internal/startup"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
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
	svcCtx     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/vid.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		//startup.Init(svcCtx)
		startup.Subscribe(svcCtx)
		logx.Infof("enabled vidsvr")
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
		vid.RegisterVidmgrMangeServer(grpcServer, vidmanage.NewVidmgrMangeServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
