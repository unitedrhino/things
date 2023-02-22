package main

import (
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
	"github.com/i-Things/things/src/rulesvr/ruledirect"

	ruleengine "github.com/i-Things/things/src/rulesvr/internal/server/ruleengine"
	scenelinkage "github.com/i-Things/things/src/rulesvr/internal/server/scenelinkage"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	svcCtx := ruledirect.GetSvcCtx()
	c := svcCtx.Config

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rule.RegisterSceneLinkageServer(grpcServer, scenelinkage.NewSceneLinkageServer(svcCtx))
		rule.RegisterRuleEngineServer(grpcServer, ruleengine.NewRuleEngineServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(errors.ErrorInterceptor)

	fmt.Printf("Starting rulesvr server at %s...\n", c.ListenOn)
	s.Start()
}
