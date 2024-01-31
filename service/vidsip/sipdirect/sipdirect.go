package sipdirect

import (
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	client "github.com/i-Things/things/service/vidsip/client/sipmanage"
	"github.com/i-Things/things/service/vidsip/internal/config"
	server "github.com/i-Things/things/service/vidsip/internal/server/sipmanage"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"
	"sync"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	SipInfo client.SipInfo
)

type Config = config.Config

var (
	c          config.Config
	svcCtx     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/sip.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		logx.Infof("enabled vidsip")
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
	fmt.Println("[---test--] vidsvr  svcCtx.Config  c.Mode:", c.Mode)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		sip.RegisterSipManageServer(grpcServer, server.NewSipManageServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
			fmt.Println("[---test--] vidsip  svcCtx.Config")
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func NewSipManage(runSvr bool) client.SipManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	vidSvr := client.NewDirectSipManage(svcCtx, server.NewSipManageServer(svcCtx))
	return vidSvr
}
