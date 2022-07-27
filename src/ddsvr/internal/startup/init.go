package startup

import (
	"context"
	"fmt"
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/event"
	"github.com/i-Things/things/src/ddsvr/internal/handler"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/devLink"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/innerLink"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"os"
)

func NewDd(c config.Config) {
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

	fmt.Printf("Starting ddSvr at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
