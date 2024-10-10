package dmdirect

import (
	"fmt"
	"gitee.com/unitedrhino/share/interceptors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/config"
	devicegroup "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicegroup"
	deviceinteract "gitee.com/unitedrhino/things/service/dmsvr/internal/server/deviceinteract"
	devicemanage "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemanage"
	devicemsg "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemsg"
	otamanage "gitee.com/unitedrhino/things/service/dmsvr/internal/server/otamanage"
	productmanage "gitee.com/unitedrhino/things/service/dmsvr/internal/server/productmanage"
	protocolmanage "gitee.com/unitedrhino/things/service/dmsvr/internal/server/protocolmanage"
	remoteconfig "gitee.com/unitedrhino/things/service/dmsvr/internal/server/remoteconfig"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/startup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
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
		utils.ConfMustLoad("etc/dm.yaml", &c)
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
		dm.RegisterOtaManageServer(grpcServer, otamanage.NewOtaManageServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
