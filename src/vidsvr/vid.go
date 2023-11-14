package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/viddirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := viddirect.GetSvcCtx()
	viddirect.Run(svcCtx)
}

//var configFile = flag.String("f", "etc/vid.yaml", "the config file")
//
//func main() {
//	flag.Parse()
//
//	var c config.Config
//	conf.MustLoad(*configFile, &c)
//	ctx := svc.NewServiceContext(c)
//
//	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
//		vid.RegisterVidmgrMangeServer(grpcServer, vidmgrmangeServer.NewVidmgrMangeServer(ctx))
//
//		if c.Mode == service.DevMode || c.Mode == service.TestMode {
//			fmt.Printf("Starting rpc serverc.Mode == service.DevMode\n")
//			reflection.Register(grpcServer)
//		}
//	})
//	defer s.Stop()
//
//	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
//	s.Start()
//}
