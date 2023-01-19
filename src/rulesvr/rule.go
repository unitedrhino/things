package main

import (
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
	"github.com/i-Things/things/src/rulesvr/ruledirect"

	flow "github.com/i-Things/things/src/rulesvr/internal/server/flow"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	svcCtx := ruledirect.GetCtxSvc()
	c := svcCtx.Config

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rule.RegisterFlowServer(grpcServer, flow.NewFlowServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
