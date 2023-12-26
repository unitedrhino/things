package sysdirect

import (
	"fmt"
	"github.com/i-Things/things/src/syssvr/internal/config"
	appmanageServer "github.com/i-Things/things/src/syssvr/internal/server/appmanage"
	areamanageServer "github.com/i-Things/things/src/syssvr/internal/server/areamanage"
	commonServer "github.com/i-Things/things/src/syssvr/internal/server/common"
	logServer "github.com/i-Things/things/src/syssvr/internal/server/log"
	modulemanageServer "github.com/i-Things/things/src/syssvr/internal/server/modulemanage"
	projectmanageServer "github.com/i-Things/things/src/syssvr/internal/server/projectmanage"
	rolemanageServer "github.com/i-Things/things/src/syssvr/internal/server/rolemanage"
	tenantmanageServer "github.com/i-Things/things/src/syssvr/internal/server/tenantmanage"
	usermanageServer "github.com/i-Things/things/src/syssvr/internal/server/usermanage"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
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
		conf.MustLoad("etc/sys.yaml", &c)
		ctxSvc = svc.NewServiceContext(c)
		logx.Infof("enabled syssvr")
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
		sys.RegisterUserManageServer(grpcServer, usermanageServer.NewUserManageServer(svcCtx))
		sys.RegisterRoleManageServer(grpcServer, rolemanageServer.NewRoleManageServer(svcCtx))
		sys.RegisterAppManageServer(grpcServer, appmanageServer.NewAppManageServer(svcCtx))
		sys.RegisterCommonServer(grpcServer, commonServer.NewCommonServer(svcCtx))
		sys.RegisterLogServer(grpcServer, logServer.NewLogServer(svcCtx))
		sys.RegisterModuleManageServer(grpcServer, modulemanageServer.NewModuleManageServer(svcCtx))
		sys.RegisterProjectManageServer(grpcServer, projectmanageServer.NewProjectManageServer(svcCtx))
		sys.RegisterAreaManageServer(grpcServer, areamanageServer.NewAreaManageServer(svcCtx))
		sys.RegisterTenantManageServer(grpcServer, tenantmanageServer.NewTenantManageServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
