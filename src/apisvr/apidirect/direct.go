package apidirect

import (
	"fmt"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/handler"
	"github.com/i-Things/things/src/apisvr/internal/handler/frontend"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/ddsvr/dddirect"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

type (
	Config = config.Config
	ApiCtx struct {
		Server *rest.Server
	}
)

var (
	c config.Config
)

func NewApi(apiCtx ApiCtx) *rest.Server {
	conf.MustLoad("etc/api.yaml", &c)
	if c.DdEnable == true {
		go runDdSvr()
	}
	return runApi(apiCtx)
}

func runApi(apiCtx ApiCtx) *rest.Server {
	ctx := svc.NewServiceContext(c)
	var server = apiCtx.Server
	if server == nil {
		server = rest.MustNewServer(c.RestConf, rest.WithCors("*"),
			rest.WithNotFoundHandler(frontend.FrontendHandler(ctx)),
		)
	}
	defer server.Stop()
	server.Use(ctx.Record)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n", c.Host, c.Port)

	return server
}
func runDdSvr() {
	dddirect.NewDd()
}
