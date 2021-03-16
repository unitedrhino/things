package main

import (
	"flag"
	"fmt"
	"yl/shared/utils"
	"yl/user/common"

	"yl/user/api/internal/config"
	"yl/user/api/internal/handler"
	"yl/user/api/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	common.UserID = utils.NewSnowFlake(c.NodeID)
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
