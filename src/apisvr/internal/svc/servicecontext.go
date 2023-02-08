package svc

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/verify"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	deviceinteract "github.com/i-Things/things/src/disvr/client/deviceinteract"
	devicemsg "github.com/i-Things/things/src/disvr/client/devicemsg"
	"github.com/i-Things/things/src/disvr/didirect"
	deviceauth "github.com/i-Things/things/src/dmsvr/client/deviceauth"
	devicegroup "github.com/i-Things/things/src/dmsvr/client/devicegroup"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	remoteconfig "github.com/i-Things/things/src/dmsvr/client/remoteconfig"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	api "github.com/i-Things/things/src/syssvr/client/api"
	common "github.com/i-Things/things/src/syssvr/client/common"
	log "github.com/i-Things/things/src/syssvr/client/log"
	menu "github.com/i-Things/things/src/syssvr/client/menu"
	role "github.com/i-Things/things/src/syssvr/client/role"
	user "github.com/i-Things/things/src/syssvr/client/user"
	"github.com/i-Things/things/src/syssvr/sysdirect"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config         config.Config
	CheckToken     rest.Middleware
	DmManage       rest.Middleware
	UserRpc        user.User
	RoleRpc        role.Role
	MenuRpc        menu.Menu
	DeviceM        devicemanage.DeviceManage
	DeviceA        deviceauth.DeviceAuth
	ProductM       productmanage.ProductManage
	DeviceMsg      devicemsg.DeviceMsg
	DeviceInteract deviceinteract.DeviceInteract
	Captcha        *verify.Captcha
	OSS            oss.OSSer
	DeviceG        devicegroup.DeviceGroup
	RemoteConfig   remoteconfig.RemoteConfig
	Common         common.Common
	LogRpc         log.Log
	ApiRpc         api.Api
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		deviceM        devicemanage.DeviceManage
		productM       productmanage.ProductManage
		deviceA        deviceauth.DeviceAuth
		deviceMsg      devicemsg.DeviceMsg
		deviceInteract deviceinteract.DeviceInteract
		deviceG        devicegroup.DeviceGroup
		remoteConfig   remoteconfig.RemoteConfig
		sysCommon      common.Common
	)
	var ur user.User
	var ro role.Role
	var me menu.Menu
	var lo log.Log
	var ap api.Api
	//var me menu.Menu
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc {
			deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
			productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceA = deviceauth.NewDeviceAuth(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))
			remoteConfig = remoteconfig.NewRemoteConfig(zrpc.MustNewClient(c.DmRpc.Conf))
			sysCommon = common.NewCommon(zrpc.MustNewClient(c.DmRpc.Conf))
		} else {
			deviceM = dmdirect.NewDeviceManage()
			productM = dmdirect.NewProductManage()
			deviceA = dmdirect.NewDeviceAuth()
			deviceG = dmdirect.NewDeviceGroup()
			remoteConfig = dmdirect.NewRemoteConfig()
			sysCommon = sysdirect.NewCommon()
		}
	}
	if c.SysRpc.Enable {
		if c.SysRpc.Mode == conf.ClientModeGrpc {
			ur = user.NewUser(zrpc.MustNewClient(c.SysRpc.Conf))
			ro = role.NewRole(zrpc.MustNewClient(c.SysRpc.Conf))
			me = menu.NewMenu(zrpc.MustNewClient(c.SysRpc.Conf))
			lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
			ap = api.NewApi(zrpc.MustNewClient(c.SysRpc.Conf))
		} else {
			ur = sysdirect.NewUser()
			ro = sysdirect.NewRole()
			me = sysdirect.NewMenu()
			lo = sysdirect.NewLog()
			ap = sysdirect.NewApi()
		}
	}
	if c.DiRpc.Enable {
		if c.DiRpc.Mode == conf.ClientModeGrpc {
			deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DiRpc.Conf))
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DiRpc.Conf))

		} else {
			deviceMsg = didirect.NewDeviceMsg()
			deviceInteract = didirect.NewDeviceInteract()
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
	jwt.TimeFunc = func() time.Time {
		return time.Now()
	}
	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:         c,
		CheckToken:     middleware.NewCheckTokenMiddleware(c, ur, lo).Handle,
		UserRpc:        ur,
		RoleRpc:        ro,
		MenuRpc:        me,
		ProductM:       productM,
		DeviceM:        deviceM,
		Captcha:        captcha,
		DeviceInteract: deviceInteract,
		DeviceMsg:      deviceMsg,
		DeviceA:        deviceA,
		DeviceG:        deviceG,
		RemoteConfig:   remoteConfig,
		Common:         sysCommon,
		LogRpc:         lo,
		ApiRpc:         ap,
		//OSS:        ossClient,
	}
}
