package svc

import (
	"context"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/verify"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/usersvr/user"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

type ServiceContext struct {
	Config     config.Config
	CheckToken rest.Middleware
	Record     rest.Middleware
	DmManage   rest.Middleware
	UserRpc    user.User
	DmRpc      dm.Dm
	DcRpc      dc.Dc
	Captcha    *verify.Captcha
	OSS        oss.OSSer
}

func NewServiceContext(c config.Config) *ServiceContext {
	var dr dm.Dm
	var ur user.User
	var dcSvr dc.Dc
	if c.DmRpc.Enable {
		dr = dm.NewDm(zrpc.MustNewClient(c.DmRpc.Conf))
	}
	if c.DcRpc.Enable {
		dcSvr = dc.NewDc(zrpc.MustNewClient(c.DcRpc.Conf))
	}
	if c.UserRpc.Enable {
		ur = user.NewUser(zrpc.MustNewClient(c.UserRpc.Conf))
	}
	ossClient, err := oss.NewOss(c.OSS)
	if err != nil {
		logx.Error("NewOss err", err)
		os.Exit(-1)
	}
	oss.InitBuckets(context.TODO(), ossClient)
	if err != nil {
		logx.Error("InitBuckets err", err)
		os.Exit(-1)
	}
	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:     c,
		CheckToken: middleware.NewCheckTokenMiddleware(ur).Handle,
		Record:     middleware.NewRecordMiddleware().Handle,
		DmManage:   middleware.NewDmManageMiddleware().Handle,
		UserRpc:    ur,
		DmRpc:      dr,
		Captcha:    captcha,
		DcRpc:      dcSvr,
		OSS:        ossClient,
	}
}
