package svc

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/verify"
	ws "github.com/i-Things/things/shared/websocket"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	"github.com/i-Things/things/src/disvr/client/deviceinteract"
	"github.com/i-Things/things/src/disvr/client/devicemsg"
	"github.com/i-Things/things/src/disvr/didirect"
	"github.com/i-Things/things/src/dmsvr/client/deviceauth"
	alarmcenter "github.com/i-Things/things/src/rulesvr/client/alarmcenter"
	scenelinkage "github.com/i-Things/things/src/rulesvr/client/scenelinkage"
	"github.com/i-Things/things/src/rulesvr/ruledirect"
	"github.com/zeromicro/go-zero/core/logx"
	"os"

	"github.com/i-Things/things/src/dmsvr/client/devicegroup"
	"github.com/i-Things/things/src/dmsvr/client/devicemanage"
	"github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/client/remoteconfig"
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

func init() {
	jwt.TimeFunc = func() time.Time {
		return time.Now()
	}
}

type SvrClient struct {
	UserRpc        user.User
	RoleRpc        role.Role
	MenuRpc        menu.Menu
	DeviceM        devicemanage.DeviceManage
	DeviceA        deviceauth.DeviceAuth
	ProductM       productmanage.ProductManage
	DeviceMsg      devicemsg.DeviceMsg
	DeviceInteract deviceinteract.DeviceInteract
	DeviceG        devicegroup.DeviceGroup
	RemoteConfig   remoteconfig.RemoteConfig
	Common         common.Common
	LogRpc         log.Log
	ApiRpc         api.Api
	Scene          scenelinkage.SceneLinkage
	Alarm          alarmcenter.AlarmCenter
}

type ServiceContext struct {
	SvrClient
	Ws             *ws.Server
	Config         config.Config
	SetupWare      rest.Middleware
	CheckTokenWare rest.Middleware
	TeardownWare   rest.Middleware
	Captcha        *verify.Captcha
	OssClient      *oss.Client
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
		scene          scenelinkage.SceneLinkage
		alarm          alarmcenter.AlarmCenter
	)
	var ur user.User
	var ro role.Role
	var me menu.Menu
	var lo log.Log
	var ap api.Api
	ws.StartWsDp(false)
	//var me menu.Menu
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc { //服务模式
			productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceA = deviceauth.NewDeviceAuth(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))
			remoteConfig = remoteconfig.NewRemoteConfig(zrpc.MustNewClient(c.DmRpc.Conf))
		} else {
			deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
			productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
			deviceA = dmdirect.NewDeviceAuth(c.DmRpc.RunProxy)
			deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
			remoteConfig = dmdirect.NewRemoteConfig(c.DmRpc.RunProxy)
		}
	}
	if c.RuleRpc.Enable {
		if c.RuleRpc.Mode == conf.ClientModeGrpc {
			scene = scenelinkage.NewSceneLinkage(zrpc.MustNewClient(c.RuleRpc.Conf))
			alarm = alarmcenter.NewAlarmCenter(zrpc.MustNewClient(c.RuleRpc.Conf))
		} else {
			scene = ruledirect.NewSceneLinkage(c.RuleRpc.RunProxy)
			alarm = ruledirect.NewAlarmCenter(c.RuleRpc.RunProxy)
		}
	}
	if c.SysRpc.Enable {
		if c.SysRpc.Mode == conf.ClientModeGrpc {
			ur = user.NewUser(zrpc.MustNewClient(c.SysRpc.Conf))
			ro = role.NewRole(zrpc.MustNewClient(c.SysRpc.Conf))
			me = menu.NewMenu(zrpc.MustNewClient(c.SysRpc.Conf))
			lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
			ap = api.NewApi(zrpc.MustNewClient(c.SysRpc.Conf))
			sysCommon = common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf))
		} else {
			ur = sysdirect.NewUser(c.SysRpc.RunProxy)
			ro = sysdirect.NewRole(c.SysRpc.RunProxy)
			me = sysdirect.NewMenu(c.SysRpc.RunProxy)
			lo = sysdirect.NewLog(c.SysRpc.RunProxy)
			ap = sysdirect.NewApi(c.SysRpc.RunProxy)
			sysCommon = sysdirect.NewCommon(c.SysRpc.RunProxy)
		}
	}
	if c.DiRpc.Enable {
		if c.DiRpc.Mode == conf.ClientModeGrpc {
			deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DiRpc.Conf))
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DiRpc.Conf))

		} else {
			deviceMsg = didirect.NewDeviceMsg(c.DiRpc.RunProxy)
			deviceInteract = didirect.NewDeviceInteract(c.DiRpc.RunProxy)
		}
	}

	ossClient := oss.NewOssClient(c.OssConf)
	if ossClient == nil {
		logx.Error("NewOss err")
		os.Exit(-1)
	}

	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:         c,
		SetupWare:      middleware.NewSetupWareMiddleware(c, lo).Handle,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware(c, ur, ro).Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware(c, lo).Handle,
		Captcha:        captcha,
		OssClient:      ossClient,
		SvrClient: SvrClient{
			UserRpc:        ur,
			RoleRpc:        ro,
			MenuRpc:        me,
			ProductM:       productM,
			DeviceM:        deviceM,
			DeviceInteract: deviceInteract,
			DeviceMsg:      deviceMsg,
			DeviceA:        deviceA,
			DeviceG:        deviceG,
			RemoteConfig:   remoteConfig,
			Common:         sysCommon,
			Scene:          scene,
			Alarm:          alarm,
			LogRpc:         lo,
			ApiRpc:         ap,
		},
		//OSS:        ossClient,
	}
}
