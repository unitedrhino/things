package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/cast"
	"yl/shared/errors"
	"yl/src/usersvr/internal/config"
	"yl/src/usersvr/internal/server"
	"yl/src/usersvr/internal/svc"
	"yl/src/usersvr/user"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/usersvr.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	conf,_ := json.Marshal(c)
	fmt.Printf("config:%s\n",cast.ToString(conf))
	ctx := svc.NewServiceContext(c)
	srv := server.NewUserServer(ctx)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, srv)
	})
	s.AddUnaryInterceptors(errors.ErrorInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
