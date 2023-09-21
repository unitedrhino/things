package apidirect

import (
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/handler"
	"github.com/i-Things/things/src/apisvr/internal/handler/system/proxy"
	"github.com/i-Things/things/src/apisvr/internal/startup"
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
		SvcCtx *ServiceContext
	}
)

var (
	c config.Config
)

func NewApi(apiCtx ApiCtx) ApiCtx {
	conf.MustLoad("etc/api.yaml", &c)
	apiCtx = runApi(apiCtx)
	if c.DdEnable == true {
		go runDdSvr()
	}
	return apiCtx
}

func runApi(apiCtx ApiCtx) ApiCtx {
	var server = apiCtx.Server
	ctx := svc.NewServiceContext(c)
	apiCtx.SvcCtx = ctx
	if server == nil {
		server = rest.MustNewServer(c.RestConf, rest.WithCors("*"),
			rest.WithNotFoundHandler(proxy.Handler(ctx)),
		)
		apiCtx.Server = server
	}
	handler.RegisterHandlers(server, ctx)
	//ota附件处理
	startup.StartOtaChanWalk(apiCtx.SvcCtx)
	return apiCtx
}
func runDdSvr() {
	dddirect.NewDd()
}
