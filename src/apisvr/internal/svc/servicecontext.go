package svc

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/verify"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/i-Things/things/src/dcsvr/dcdirect"
	"github.com/i-Things/things/src/dmsvr/dmclient"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/usersvr/userclient"
	"github.com/i-Things/things/src/usersvr/userdirect"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config     config.Configs
	CheckToken rest.Middleware
	Record     rest.Middleware
	DmManage   rest.Middleware
	UserRpc    userclient.User
	DmRpc      dmclient.Dm
	DcRpc      dc.Dc
	Captcha    *verify.Captcha
	OSS        oss.OSSer
}

func NewServiceContext(c config.Configs) *ServiceContext {
	var dr dmclient.Dm
	var ur userclient.User
	var dcSvr dc.Dc
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc {
			dr = dmclient.NewDm(zrpc.MustNewClient(c.DmRpc.Conf))
		} else {
			dr = dmdirect.NewDm(&c.DmSvr)
		}
	}
	if c.UserRpc.Enable {
		if c.UserRpc.Mode == conf.ClientModeGrpc {
			ur = userclient.NewUser(zrpc.MustNewClient(c.UserRpc.Conf))
		} else {
			ur = userdirect.NewUser(c.UserSvr)
		}
	}
	if c.DcRpc.Enable {
		if c.DcSvr.Mode == conf.ClientModeGrpc {
			dcSvr = dc.NewDc(zrpc.MustNewClient(c.DcRpc.Conf))
		} else {
			dcSvr = dcdirect.NewDc(c.DcSvr)
		}
	}
	//ossClient, err := oss.NewOss(c.OSS)
	//if err != nil {
	//	logx.Error("NewOss err", err)
	//	os.Exit(-1)
	//}
	//oss.InitBuckets(context.TODO(), ossClient)
	//if err != nil {
	//	logx.Error("InitBuckets err", err)
	//	os.Exit(-1)
	//}
	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:     c,
		CheckToken: middleware.NewCheckTokenMiddleware(ur).Handle,
		Record:     middleware.NewRecordMiddleware().Handle,
		UserRpc:    ur,
		DmRpc:      dr,
		Captcha:    captcha,
		DcRpc:      dcSvr,
		//OSS:        ossClient,
	}
}
