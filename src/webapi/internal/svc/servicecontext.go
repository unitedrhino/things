package svc

import (
	"gitee.com/godLei6/things/shared/verify"
	"gitee.com/godLei6/things/src/dcsvr/dcclient"
	"gitee.com/godLei6/things/src/dmsvr/dmclient"
	"gitee.com/godLei6/things/src/usersvr/userclient"
	"gitee.com/godLei6/things/src/webapi/internal/config"
	"gitee.com/godLei6/things/src/webapi/internal/middleware"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
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
