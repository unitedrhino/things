package apidirect

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	ws "github.com/i-Things/things/shared/websocket"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/handler"
	"github.com/i-Things/things/src/apisvr/internal/handler/system/proxy"
	"github.com/i-Things/things/src/apisvr/internal/repo/event/appDeviceEvent"
	"github.com/i-Things/things/src/apisvr/internal/repo/event/subApp"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/ddsvr/dddirect"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"log"
	"os"
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
	if apiCtx.SvcCtx.Ws == nil {
		apiCtx.SvcCtx.Ws = ws.MustNewServer(c.RestConf)
	}
	handler.RegisterHandlers(server, ctx)
	handler.RegisterWsHandlers(apiCtx.SvcCtx.Ws, ctx)
	subAppCli, err := subApp.NewSubApp(ctx.Config.Event)
	if err != nil {
		logx.Error("NewSubApp err", err)
		os.Exit(-1)
	}
	err = subAppCli.Subscribe(func(ctx1 context.Context) subApp.AppSubEvent {
		return appDeviceEvent.NewAppDeviceHandle(ctx1, ctx)
	})
	if err != nil {
		log.Fatalf("%v.subApp.Subscribe err:%v",
			utils.FuncName(), err)
	}
	return apiCtx
}
func runDdSvr() {
	dddirect.NewDd()
}
