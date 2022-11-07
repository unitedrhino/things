package apidirect

import (
	"fmt"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/handler"
	"github.com/i-Things/things/src/apisvr/internal/handler/system/proxy"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/ddsvr/dddirect"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

type (
	Config         = config.Config
	ServiceContext = svc.ServiceContext
	ApiCtx         struct {
		Server *rest.Server
		Svc    *ServiceContext
	}
)

var (
	c config.Config
)

func NewApi(apiCtx ApiCtx) ApiCtx {
	conf.MustLoad("etc/api.yaml", &c)
	if c.DdEnable == true {
		go runDdSvr()
	}
	return runApi(apiCtx)
}

func runApi(apiCtx ApiCtx) ApiCtx {
	var server = apiCtx.Server
	ctx := svc.NewServiceContext(c)
	apiCtx.Svc = ctx
	if server == nil {
		server = rest.MustNewServer(c.RestConf, rest.WithCors("*"),
			rest.WithNotFoundHandler(proxy.Handler(ctx)),
		)
		apiCtx.Server = server
	}
	defer server.Stop()
	server.Use(ctx.Record)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n", c.Host, c.Port)

	return apiCtx
}
func runDdSvr() {
	dddirect.NewDd()
}
