package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/i-Things/things/src/ddsvr/internal/event"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/devLink"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/innerLink"
	"github.com/zeromicro/go-zero/core/logx"
	"os"

	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/handler"
	"github.com/i-Things/things/src/ddsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dd.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcCtx := svc.NewServiceContext(c)
	err := svcCtx.DevLink.SubScribe(func(ctx context.Context) devLink.DevSubHandle {
		return event.NewDeviceSubServer(svcCtx, ctx)
	})
	if err != nil {
		logx.Error("DevLink Subscribe failure", err)
		os.Exit(-1)
	}
	err = svcCtx.InnerLink.Subscribe(func(ctx context.Context) innerLink.InnerSubHandle {
		return event.NewInnerSubServer(svcCtx, ctx)
	})
	if err != nil {
		logx.Error("InnerLink Subscribe failure", err)
		os.Exit(-1)
	}
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, svcCtx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
