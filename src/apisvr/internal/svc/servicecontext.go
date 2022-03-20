package svc

import (
	"github.com/i-Things/things/shared/verify"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	"github.com/i-Things/things/src/dcsvr/dcclient"
	"github.com/i-Things/things/src/dmsvr/dmclient"
	"github.com/i-Things/things/src/usersvr/userclient"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config     config.Config
	CheckToken rest.Middleware
	Record     rest.Middleware
	DmManage   rest.Middleware
	UserRpc    userclient.User
	DmRpc      dmclient.Dm
	DcRpc      dcclient.Dc
	Captcha    *verify.Captcha
}

func NewServiceContext(c config.Config) *ServiceContext {
	var dr dmclient.Dm
	var ur userclient.User
	var dc dcclient.Dc
	if c.DmRpc.Enable {
		dr = dmclient.NewDm(zrpc.MustNewClient(c.DmRpc.Conf))
	}
	if c.DcRpc.Enable {
		dc = dcclient.NewDc(zrpc.MustNewClient(c.DcRpc.Conf))
	}
	if c.UserRpc.Enable {
		ur = userclient.NewUser(zrpc.MustNewClient(c.UserRpc.Conf))
	}
	captcha := verify.NewCaptcha(c.ImgHeight, c.ImgWidth, c.KeyLong, c.CacheRedis, time.Duration(c.KeepTime)*time.Second)
	return &ServiceContext{
		Config:     c,
		CheckToken: middleware.NewCheckTokenMiddleware(ur).Handle,
		Record:     middleware.NewRecordMiddleware().Handle,
		DmManage:   middleware.NewDmManageMiddleware().Handle,
		UserRpc:    ur,
		DmRpc:      dr,
		Captcha:    captcha,
		DcRpc:      dc,
	}
}
