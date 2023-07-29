package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/stocksvr/stockdirect"
)

//var configFile = flag.String("f", "etc/stock.yaml", "the config file")

func main() {
	//flag.Parse()
	//
	//var c config.Config
	//conf.MustLoad(*configFile, &c)
	//ctx := svc.NewServiceContext(c)
	//
	//s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
	//	stock.RegisterStockServer(grpcServer, server.NewStockServer(ctx))
	//
	//	if c.Mode == service.DevMode || c.Mode == service.TestMode {
	//		reflection.Register(grpcServer)
	//	}
	//})
	//defer s.Stop()
	//
	//fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	//s.Start()
	defer utils.Recover(context.Background())
	svcCtx := stockdirect.GetSvcCtx()
	stockdirect.RunServer(svcCtx)
}
