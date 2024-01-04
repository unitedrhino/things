package viddirect

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/config"
	mgrconfig "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrconfigmanage"
	mgrgbsip "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrgbsipmanage"
	mgrinfo "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrinfomanage"
	mgrstream "github.com/i-Things/things/src/vidsvr/internal/server/vidmgrstreammanage"
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
	c          config.Config
	svcCtx     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/vid.yaml", &c)
		svcCtx = svc.NewServiceContext(c)

		startup.Subscribe(svcCtx)
		logx.Infof("enabled vidsvr")
	})
	return svcCtx
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		utils.Go(context.Background(), ApiRun) //golang 后台执行
		go Run(svcCtx)
	})
}

func Run(svcCtx *svc.ServiceContext) {
	c := svcCtx.Config
	fmt.Println("[---test--] vidsvr  svcCtx.Config  c.Mode:", c.Mode)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		vid.RegisterVidmgrInfoManageServer(grpcServer, mgrinfo.NewVidmgrInfoManageServer(svcCtx))
		vid.RegisterVidmgrConfigManageServer(grpcServer, mgrconfig.NewVidmgrConfigManageServer(svcCtx))
		vid.RegisterVidmgrStreamManageServer(grpcServer, mgrstream.NewVidmgrStreamManageServer(svcCtx))
		vid.RegisterVidmgrGbsipManageServer(grpcServer, mgrgbsip.NewVidmgrGbsipManageServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
			fmt.Println("[---test--] vidsvr  svcCtx.Config")
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
