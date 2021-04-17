package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/cast"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/common"

	"yl/src/user/rpc/internal/config"
	"yl/src/user/rpc/internal/server"
	"yl/src/user/rpc/internal/svc"
	"yl/src/user/rpc/user"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	conf,_ := json.Marshal(c)
	fmt.Printf("config:%s\n",cast.ToString(conf))
	ctx := svc.NewServiceContext(c)
	srv := server.NewUserServer(ctx)
	common.UserID = utils.NewSnowFlake(c.NodeID)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, srv)
	})
	s.AddUnaryInterceptors(errors.ErrorInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
