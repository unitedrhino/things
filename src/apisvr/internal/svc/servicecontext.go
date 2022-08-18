package svc

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/verify"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	deviceinteract "github.com/i-Things/things/src/disvr/client/deviceinteract"
	devicemsg "github.com/i-Things/things/src/disvr/client/devicemsg"
	"github.com/i-Things/things/src/disvr/didirect"
	deviceauth "github.com/i-Things/things/src/dmsvr/client/deviceauth"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/usersvr/userclient"
	"github.com/i-Things/things/src/usersvr/userdirect"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config         config.Configs
	CheckToken     rest.Middleware
	Record         rest.Middleware
	DmManage       rest.Middleware
	UserRpc        userclient.User
	DeviceM        devicemanage.DeviceManage
	DeviceA        deviceauth.DeviceAuth
	ProductM       productmanage.ProductManage
	DeviceMsg      devicemsg.DeviceMsg
	DeviceInteract deviceinteract.DeviceInteract
	Captcha        *verify.Captcha
	OSS            oss.OSSer
}

func NewServiceContext(c config.Configs) *ServiceContext {
	var (
		deviceM        devicemanage.DeviceManage
		productM       productmanage.ProductManage
		deviceA        deviceauth.DeviceAuth
		deviceMsg      devicemsg.DeviceMsg
		deviceInteract deviceinteract.DeviceInteract
	)
	var ur userclient.User
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc {
			deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
			productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceA = deviceauth.NewDeviceAuth(zrpc.MustNewClient(c.DmRpc.Conf))
		} else {
			deviceM = dmdirect.NewDeviceManage(c.DmSvr)
			productM = dmdirect.NewProductManage(c.DmSvr)
			deviceA = dmdirect.NewDeviceAuth(c.DmSvr)
		}
	}
	if c.UserRpc.Enable {
		if c.UserRpc.Mode == conf.ClientModeGrpc {
			ur = userclient.NewUser(zrpc.MustNewClient(c.UserRpc.Conf))
		} else {
			ur = userdirect.NewUser(c.UserSvr)
		}
	}
	if c.DiRpc.Enable {
		if c.DiRpc.Mode == conf.ClientModeGrpc {
			deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DiRpc.Conf))
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DiRpc.Conf))

		} else {
			deviceMsg = didirect.NewDeviceMsg(c.DiSvr)
			deviceInteract = didirect.NewDeviceInteract(c.DiSvr)
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
		Config:         c,
		CheckToken:     middleware.NewCheckTokenMiddleware(ur).Handle,
		Record:         middleware.NewRecordMiddleware().Handle,
		UserRpc:        ur,
		ProductM:       productM,
		DeviceM:        deviceM,
		Captcha:        captcha,
		DeviceInteract: deviceInteract,
		DeviceMsg:      deviceMsg,
		DeviceA:        deviceA,
		//OSS:        ossClient,
	}
}
