package svc

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/verify"
	ws "github.com/i-Things/things/shared/websocket"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	"github.com/i-Things/things/src/dgsvr/client/deviceauth"
	"github.com/i-Things/things/src/dgsvr/dgdirect"
	"github.com/i-Things/things/src/dmsvr/client/devicegroup"
	"github.com/i-Things/things/src/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/src/dmsvr/client/devicemanage"
	"github.com/i-Things/things/src/dmsvr/client/devicemsg"
	firmwaremanage "github.com/i-Things/things/src/dmsvr/client/firmwaremanage"
	otataskmanage "github.com/i-Things/things/src/dmsvr/client/otataskmanage"
	"github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/client/protocolmanage"
	"github.com/i-Things/things/src/dmsvr/client/remoteconfig"
	"github.com/i-Things/things/src/dmsvr/client/schemamanage"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/rulesvr/client/alarmcenter"
	"github.com/i-Things/things/src/rulesvr/client/scenelinkage"
	accesssManage "github.com/i-Things/things/src/syssvr/client/accessmanage"
	app "github.com/i-Things/things/src/syssvr/client/appmanage"
	"github.com/i-Things/things/src/syssvr/client/areamanage"
	common "github.com/i-Things/things/src/syssvr/client/common"
	log "github.com/i-Things/things/src/syssvr/client/log"
	module "github.com/i-Things/things/src/syssvr/client/modulemanage"
	"github.com/i-Things/things/src/syssvr/client/projectmanage"
	role "github.com/i-Things/things/src/syssvr/client/rolemanage"
	tenant "github.com/i-Things/things/src/syssvr/client/tenantmanage"
	user "github.com/i-Things/things/src/syssvr/client/usermanage"
	"github.com/i-Things/things/src/syssvr/sysdirect"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
	"github.com/i-Things/things/src/timed/timedjobsvr/timedjobdirect"
	"github.com/i-Things/things/src/timed/timedschedulersvr/client/timedscheduler"
	"github.com/i-Things/things/src/timed/timedschedulersvr/timedschedulerdirect"
	"github.com/i-Things/things/src/udsvr/client/ops"
	"github.com/i-Things/things/src/udsvr/client/rule"
	"github.com/i-Things/things/src/udsvr/uddirect"
	"github.com/i-Things/things/src/vidsvr/client/vidmgrconfigmanage"
	"github.com/i-Things/things/src/vidsvr/client/vidmgrinfomanage"
	"github.com/i-Things/things/src/vidsvr/client/vidmgrstreammanage"
	"github.com/i-Things/things/src/vidsvr/viddirect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

func init() {
	jwt.TimeFunc = func() time.Time {
		return time.Now()
	}
}

type SvrClient struct {
	TenantRpc tenant.TenantManage
	UserRpc   user.UserManage
	RoleRpc   role.RoleManage
	AccessRpc accesssManage.AccessManage
	AppRpc    app.AppManage
	ModuleRpc module.ModuleManage
	LogRpc    log.Log
	VidmgrM   vidmgrinfomanage.VidmgrInfoManage
	VidmgrC   vidmgrconfigmanage.VidmgrConfigManage
	VidmgrS   vidmgrstreammanage.VidmgrStreamManage

	ProjectM  projectmanage.ProjectManage
	ProtocolM protocolmanage.ProtocolManage
	AreaM     areamanage.AreaManage
	ProductM  productmanage.ProductManage
	SchemaM   schemamanage.SchemaManage
	DeviceM   devicemanage.DeviceManage
	DeviceA   deviceauth.DeviceAuth
	DeviceG   devicegroup.DeviceGroup

	DeviceMsg      devicemsg.DeviceMsg
	DeviceInteract deviceinteract.DeviceInteract

	RemoteConfig remoteconfig.RemoteConfig
	Common       common.Common

	Rule           rule.Rule
	Ops            ops.Ops
	Scene          scenelinkage.SceneLinkage
	Alarm          alarmcenter.AlarmCenter
	Timedscheduler timedscheduler.Timedscheduler
	TimedJob       timedmanage.TimedManage
}

type ServiceContext struct {
	SvrClient
	Ws             *ws.Server
	Config         config.Config
	SetupWare      rest.Middleware
	CheckTokenWare rest.Middleware
	DataAuthWare   rest.Middleware
	TeardownWare   rest.Middleware
	CheckApiWare   rest.Middleware
	Captcha        *verify.Captcha
	OssClient      *oss.Client
	FirmwareM      firmwaremanage.FirmwareManage
	OtaTaskM       otataskmanage.OtaTaskManage
	FileChan       chan int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		vidmgrM        vidmgrinfomanage.VidmgrInfoManage
		vidmgrC        vidmgrconfigmanage.VidmgrConfigManage
		vidmgrS        vidmgrstreammanage.VidmgrStreamManage
		Ops            ops.Ops
		protocolM      protocolmanage.ProtocolManage
		schemaM        schemamanage.SchemaManage
		projectM       projectmanage.ProjectManage
		areaM          areamanage.AreaManage
		productM       productmanage.ProductManage
		deviceM        devicemanage.DeviceManage
		deviceA        deviceauth.DeviceAuth
		deviceG        devicegroup.DeviceGroup
		appRpc         app.AppManage
		deviceMsg      devicemsg.DeviceMsg
		deviceInteract deviceinteract.DeviceInteract
		remoteConfig   remoteconfig.RemoteConfig
		sysCommon      common.Common
		firmwareM      firmwaremanage.FirmwareManage
		otaTaskM       otataskmanage.OtaTaskManage
		timedSchedule  timedscheduler.Timedscheduler
		timedJob       timedmanage.TimedManage
		tenantM        tenant.TenantManage
		accessM        accesssManage.AccessManage
		ic             rule.Rule
	)
	var ur user.UserManage
	var ro role.RoleManage
	var me module.ModuleManage
	var lo log.Log

	caches.InitStore(c.CacheRedis)

	ws.StartWsDp(false)

	if c.SysRpc.Enable {
		if c.SysRpc.Mode == conf.ClientModeGrpc {
			projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
			areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
			ur = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
			ro = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
			me = module.NewModuleManage(zrpc.MustNewClient(c.SysRpc.Conf))
			lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
			sysCommon = common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf))
			appRpc = app.NewAppManage(zrpc.MustNewClient(c.SysRpc.Conf))
			tenantM = tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
			accessM = accesssManage.NewAccessManage(zrpc.MustNewClient(c.SysRpc.Conf))
		} else {
			projectM = sysdirect.NewProjectManage(c.SysRpc.RunProxy)
			areaM = sysdirect.NewAreaManage(c.SysRpc.RunProxy)
			ur = sysdirect.NewUser(c.SysRpc.RunProxy)
			ro = sysdirect.NewRole(c.SysRpc.RunProxy)
			me = sysdirect.NewModule(c.SysRpc.RunProxy)
			lo = sysdirect.NewLog(c.SysRpc.RunProxy)
			sysCommon = sysdirect.NewCommon(c.SysRpc.RunProxy)
			appRpc = sysdirect.NewApp(c.SysRpc.RunProxy)
			tenantM = sysdirect.NewTenantManage(c.SysRpc.RunProxy)
			accessM = sysdirect.NewAccess(c.SysRpc.RunProxy)
		}
	}

	//var me menu.Menu
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc { //服务模式
			deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DmRpc.Conf))
			productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))
			remoteConfig = remoteconfig.NewRemoteConfig(zrpc.MustNewClient(c.DmRpc.Conf))
			firmwareM = firmwaremanage.NewFirmwareManage(zrpc.MustNewClient(c.DmRpc.Conf))
			otaTaskM = otataskmanage.NewOtaTaskManage(zrpc.MustNewClient(c.DmRpc.Conf))
			protocolM = protocolmanage.NewProtocolManage(zrpc.MustNewClient(c.DmRpc.Conf))
			schemaM = schemamanage.NewSchemaManage(zrpc.MustNewClient(c.DmRpc.Conf))
		} else { //直连模式
			deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
			deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
			deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
			productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
			deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
			remoteConfig = dmdirect.NewRemoteConfig(c.DmRpc.RunProxy)
			firmwareM = dmdirect.NewFirmwareManage(c.DmRpc.RunProxy)
			otaTaskM = dmdirect.NewOtaTaskManage(c.DmRpc.RunProxy)
			protocolM = dmdirect.NewProtocolManage(c.DmRpc.RunProxy)
			schemaM = dmdirect.NewSchemaManage(c.DmRpc.RunProxy)
		}
	}
	if c.DgRpc.Enable {
		if c.DgRpc.Mode == conf.ClientModeGrpc { //服务模式
			deviceA = deviceauth.NewDeviceAuth(zrpc.MustNewClient(c.DgRpc.Conf))

		} else { //直连模式
			deviceA = dgdirect.NewDeviceAuth(c.DgRpc.RunProxy)
		}
	}
	if c.UdRpc.Enable {
		if c.UdRpc.Mode == conf.ClientModeGrpc {
			ic = rule.NewRule(zrpc.MustNewClient(c.UdRpc.Conf))
			Ops = ops.NewOps(zrpc.MustNewClient(c.UdRpc.Conf))
		} else {
			ic = uddirect.NewRule(c.UdRpc.RunProxy)
			Ops = uddirect.NewOps(c.UdRpc.RunProxy)
		}
	}

	if c.VidRpc.Enable {
		if c.VidRpc.Mode == conf.ClientModeGrpc {
			vidmgrM = vidmgrinfomanage.NewVidmgrInfoManage(zrpc.MustNewClient(c.VidRpc.Conf))
			vidmgrC = vidmgrconfigmanage.NewVidmgrConfigManage(zrpc.MustNewClient(c.VidRpc.Conf))
			vidmgrS = vidmgrstreammanage.NewVidmgrStreamManage(zrpc.MustNewClient(c.VidRpc.Conf))
		} else {
			vidmgrM = viddirect.NewVidmgrManage(c.VidRpc.RunProxy)
			vidmgrC = viddirect.NewVidmgrConfigManage(c.VidRpc.RunProxy)
			vidmgrS = viddirect.NewVidmgrStreamManage(c.VidRpc.RunProxy)
			viddirect.ApiDirectRun()
		}
	}

	if c.TimedSchedulerRpc.Enable {
		if c.TimedSchedulerRpc.Mode == conf.ClientModeGrpc {
			timedSchedule = timedscheduler.NewTimedscheduler(zrpc.MustNewClient(c.TimedSchedulerRpc.Conf))
		} else {
			timedSchedule = timedschedulerdirect.NewScheduler(c.TimedSchedulerRpc.RunProxy)
		}
	}
	if c.TimedJobRpc.Enable {
		if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
			timedJob = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
		} else {
			timedJob = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
		}
	}

	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}

	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:         c,
		SetupWare:      middleware.NewSetupWareMiddleware(c, lo).Handle,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware(c, ur, ro).Handle,
		DataAuthWare:   middleware.NewDataAuthWareMiddleware(c).Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware(c, lo).Handle,
		CheckApiWare:   middleware.NewCheckApiWareMiddleware().Handle,
		Captcha:        captcha,
		OssClient:      ossClient,
		FirmwareM:      firmwareM,
		OtaTaskM:       otaTaskM,
		Ws:             ws.MustNewServer(c.RestConf),
		SvrClient: SvrClient{
			TenantRpc:      tenantM,
			AppRpc:         appRpc,
			UserRpc:        ur,
			RoleRpc:        ro,
			AccessRpc:      accessM,
			ModuleRpc:      me,
			LogRpc:         lo,
			Timedscheduler: timedSchedule,
			TimedJob:       timedJob,
			VidmgrM:        vidmgrM,
			VidmgrC:        vidmgrC,
			VidmgrS:        vidmgrS,

			ProtocolM: protocolM,
			ProjectM:  projectM,
			SchemaM:   schemaM,
			AreaM:     areaM,
			ProductM:  productM,
			DeviceM:   deviceM,
			DeviceA:   deviceA,
			DeviceG:   deviceG,

			DeviceMsg:      deviceMsg,
			DeviceInteract: deviceInteract,
			RemoteConfig:   remoteConfig,
			Common:         sysCommon,
			Rule:           ic,
			Ops:            Ops,
		},
		//OSS:        ossClient,
	}
}
