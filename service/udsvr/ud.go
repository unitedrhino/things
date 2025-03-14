package main

import (
	"flag"
	"fmt"
	"gitee.com/unitedrhino/share/services"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/things/service/udsvr/internal/config"
	ruleServer "gitee.com/unitedrhino/things/service/udsvr/internal/server/rule"
	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	flag.Parse()

	var c config.Config
	utils.ConfMustLoad("etc/ud.yaml", &c)
	ctx := svc.NewServiceContext(c)

	s := services.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ud.RegisterRuleServer(grpcServer, ruleServer.NewRuleServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
