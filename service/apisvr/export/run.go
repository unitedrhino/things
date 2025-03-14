package export

import (
	"gitee.com/unitedrhino/share/services"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/config"
	"gitee.com/unitedrhino/things/service/apisvr/internal/handler"
	"gitee.com/unitedrhino/things/service/apisvr/internal/startup"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

type (
	Config         = config.Config
	ServiceContext = svc.ServiceContext
	ApiCtx         struct {
		Server        *rest.Server
		SvcCtx        *ServiceContext
		NotInitHandle bool
	}
)

var (
	c config.Config
)

func NewApi(apiCtx ApiCtx) ApiCtx {
	utils.ConfMustLoad("etc/api.yaml", &c)
	apiCtx = runApi(apiCtx)
	return apiCtx
}

func runApi(apiCtx ApiCtx) ApiCtx {
	var server = apiCtx.Server
	ctx := svc.NewServiceContext(c)
	apiCtx.SvcCtx = ctx
	if server == nil {
		server = rest.MustNewServer(c.RestConf)
		apiCtx.Server = server
	}
	if !apiCtx.NotInitHandle {
		handler.RegisterHandlers(server, ctx)
	}
	startup.Init(apiCtx.SvcCtx)
	services.InitApisvrs(server)
	return apiCtx
}
