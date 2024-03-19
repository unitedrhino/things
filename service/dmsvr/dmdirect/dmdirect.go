package dmdirect

import (
	"fmt"
	"gitee.com/i-Things/share/interceptors"
	"github.com/i-Things/things/service/dmsvr/internal/config"
	devicegroup "github.com/i-Things/things/service/dmsvr/internal/server/devicegroup"
	deviceinteract "github.com/i-Things/things/service/dmsvr/internal/server/deviceinteract"
	devicemanage "github.com/i-Things/things/service/dmsvr/internal/server/devicemanage"
	devicemsg "github.com/i-Things/things/service/dmsvr/internal/server/devicemsg"
	productmanage "github.com/i-Things/things/service/dmsvr/internal/server/productmanage"
	protocolmanage "github.com/i-Things/things/service/dmsvr/internal/server/protocolmanage"
	remoteconfig "github.com/i-Things/things/service/dmsvr/internal/server/remoteconfig"
	"github.com/i-Things/things/service/dmsvr/internal/startup"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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
		conf.MustLoad("etc/dm.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Init(svcCtx)
		logx.Infof("enabled dmsvr")
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
		dm.RegisterProtocolManageServer(grpcServer, protocolmanage.NewProtocolManageServer(svcCtx))
		dm.RegisterDeviceManageServer(grpcServer, devicemanage.NewDeviceManageServer(svcCtx))
		dm.RegisterProductManageServer(grpcServer, productmanage.NewProductManageServer(svcCtx))
		dm.RegisterRemoteConfigServer(grpcServer, remoteconfig.NewRemoteConfigServer(svcCtx))
		dm.RegisterDeviceGroupServer(grpcServer, devicegroup.NewDeviceGroupServer(svcCtx))
		dm.RegisterDeviceInteractServer(grpcServer, deviceinteract.NewDeviceInteractServer(svcCtx))
		dm.RegisterDeviceMsgServer(grpcServer, devicemsg.NewDeviceMsgServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
