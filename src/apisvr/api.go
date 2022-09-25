package main

import (
	"flag"
	"fmt"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/handler"
	"github.com/i-Things/things/src/apisvr/internal/handler/frontend"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/ddsvr/dddirect"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// var configFile = flag.String("f", "etc/api.yaml", "the config file")
var directFile = flag.String("f", "etc/direct.yaml", "the config file")

func main() {
	logx.DisableStat()
	flag.Parse()

	var c config.Configs
	conf.MustLoad(*directFile, &c)
	if c.DdSvr != nil {
		go runDdSvr(*c.DdSvr)
	}
	runApi(c)

}

func runApi(c config.Configs) {
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"),
		rest.WithNotFoundHandler(frontend.FrontendHandler(ctx)),
	)
	defer server.Stop()
	server.Use(ctx.Record)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
func runDdSvr(c dddirect.Config) {
	dddirect.NewDd(c)
}
