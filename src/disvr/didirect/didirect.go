package didirect

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/src/disvr/internal/config"
	deviceinteract "github.com/i-Things/things/src/disvr/internal/server/deviceinteract"
	devicemsg "github.com/i-Things/things/src/disvr/internal/server/devicemsg"
	"github.com/i-Things/things/src/disvr/internal/startup"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"
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
		flag.Parse()
		conf.MustLoad("etc/di.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Subscribe(svcCtx)
		logx.Infof("enabled disvr")
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
