package main

import (
	"flag"
	"fmt"
	"gitee.com/i-Things/share/utils"

	"github.com/i-Things/things/service/udsvr/internal/config"
	ruleServer "github.com/i-Things/things/service/udsvr/internal/server/rule"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	flag.Parse()

	var c config.Config
	utils.ConfMustLoad("etc/ud.yaml", &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ud.RegisterRuleServer(grpcServer, ruleServer.NewRuleServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
