package svc

import (
	"gitee.com/i-Things/core/service/syssvr/client/accessmanage"
	app "gitee.com/i-Things/core/service/syssvr/client/appmanage"
	"gitee.com/i-Things/core/service/syssvr/client/areamanage"
	"gitee.com/i-Things/core/service/syssvr/client/common"
	"gitee.com/i-Things/core/service/syssvr/client/datamanage"
	"gitee.com/i-Things/core/service/syssvr/client/log"
	module "gitee.com/i-Things/core/service/syssvr/client/modulemanage"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/sysdirect"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/client/timedscheduler"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/timedschedulerdirect"
	"gitee.com/i-Things/core/shared/caches"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/oss"
	"gitee.com/i-Things/core/shared/verify"
	ws "gitee.com/i-Things/core/shared/websocket"
	"github.com/dgrijalva/jwt-go"
	"github.com/i-Things/things/service/apisvr/internal/config"
	"github.com/i-Things/things/service/apisvr/internal/middleware"
	"github.com/i-Things/things/service/dgsvr/client/deviceauth"
	"github.com/i-Things/things/service/dgsvr/dgdirect"
	"github.com/i-Things/things/service/dmsvr/client/devicegroup"
	"github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/devicemsg"
	firmwaremanage "github.com/i-Things/things/service/dmsvr/client/firmwaremanage"
	"github.com/i-Things/things/service/dmsvr/client/otafirmwaremanage"
	"github.com/i-Things/things/service/dmsvr/client/otajobmanage"
	"github.com/i-Things/things/service/dmsvr/client/otamodulemanage"
	otataskmanage "github.com/i-Things/things/service/dmsvr/client/otataskmanage"
	"github.com/i-Things/things/service/dmsvr/client/otaupgradetaskmanage"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/client/protocolmanage"
	"github.com/i-Things/things/service/dmsvr/client/remoteconfig"
	"github.com/i-Things/things/service/dmsvr/client/schemamanage"
	"github.com/i-Things/things/service/dmsvr/dmdirect"
	"github.com/i-Things/things/service/rulesvr/client/alarmcenter"
	"github.com/i-Things/things/service/rulesvr/client/scenelinkage"
	"github.com/i-Things/things/service/udsvr/client/ops"
	"github.com/i-Things/things/service/udsvr/client/rule"
	"github.com/i-Things/things/service/udsvr/client/userdevice"
	"github.com/i-Things/things/service/udsvr/uddirect"
	"github.com/i-Things/things/service/vidsip/client/sipmanage"
	"github.com/i-Things/things/service/vidsip/sipdirect"
	"github.com/i-Things/things/service/vidsvr/client/vidmgrconfigmanage"
	"github.com/i-Things/things/service/vidsvr/client/vidmgrinfomanage"
	"github.com/i-Things/things/service/vidsvr/client/vidmgrstreammanage"
	"github.com/i-Things/things/service/vidsvr/viddirect"
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
	AppRpc    app.AppManage
	ModuleRpc module.ModuleManage
	LogRpc    log.Log
	SipRpc    sipmanage.SipManage
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
	AccessRpc      accessmanage.AccessManage
	DataM          datamanage.DataManage
	Ops            ops.Ops
	UserDevice     userdevice.UserDevice
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
	OtaFirmwareM   otafirmwaremanage.OTAFirmwareManage
	OtaTaskM       otataskmanage.OtaTaskManage
	FileChan       chan int64
	OtaJobM        otajobmanage.OTAJobManage
	TaskM          otaupgradetaskmanage.OTAUpgradeTaskManage
	OtaModuleM     otamodulemanage.OTAModuleManage
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		Ops     ops.Ops
		schemaM schemamanage.SchemaManage
		appRpc  app.AppManage
		vidmgrM vidmgrinfomanage.VidmgrInfoManage
		vidmgrC vidmgrconfigmanage.VidmgrConfigManage
		vidmgrS vidmgrstreammanage.VidmgrStreamManage

		vidSip sipmanage.SipManage

		protocolM protocolmanage.ProtocolManage
		projectM  projectmanage.ProjectManage
		areaM     areamanage.AreaManage
		productM  productmanage.ProductManage
		deviceM   devicemanage.DeviceManage
		deviceA   deviceauth.DeviceAuth
		deviceG   devicegroup.DeviceGroup

		deviceMsg      devicemsg.DeviceMsg
		deviceInteract deviceinteract.DeviceInteract
		remoteConfig   remoteconfig.RemoteConfig
		sysCommon      common.Common
		firmwareM      firmwaremanage.FirmwareManage
		otaFirmwareM   otafirmwaremanage.OTAFirmwareManage
		otaTaskM       otataskmanage.OtaTaskManage
		otaJobM        otajobmanage.OTAJobManage
		taskM          otaupgradetaskmanage.OTAUpgradeTaskManage
		otaModuleM     otamodulemanage.OTAModuleManage
		timedSchedule  timedscheduler.Timedscheduler
		timedJob       timedmanage.TimedManage
		tenantM        tenant.TenantManage
		UserDevice     userdevice.UserDevice
		Rule           rule.Rule
		DataM          datamanage.DataManage
		accessM        accessmanage.AccessManage
	)
	var ur user.UserManage
	var ro role.RoleManage
	var me module.ModuleManage
	var lo log.Log

	caches.InitStore(c.CacheRedis)

	ws.StartWsDp(false)
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
			otaFirmwareM = otafirmwaremanage.NewOTAFirmwareManage(zrpc.MustNewClient(c.DmRpc.Conf))
			otaTaskM = otataskmanage.NewOtaTaskManage(zrpc.MustNewClient(c.DmRpc.Conf))
			protocolM = protocolmanage.NewProtocolManage(zrpc.MustNewClient(c.DmRpc.Conf))
			otaJobM = otajobmanage.NewOTAJobManage(zrpc.MustNewClient(c.DmRpc.Conf))
			taskM = otaupgradetaskmanage.NewOTAUpgradeTaskManage(zrpc.MustNewClient(c.DmRpc.Conf))
			otaModuleM = otamodulemanage.NewOTAModuleManage(zrpc.MustNewClient(c.DmRpc.Conf))
			schemaM = schemamanage.NewSchemaManage(zrpc.MustNewClient(c.DmRpc.Conf))
		} else { //直连模式
			deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
			deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
			deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
			productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
			deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
			remoteConfig = dmdirect.NewRemoteConfig(c.DmRpc.RunProxy)
			firmwareM = dmdirect.NewFirmwareManage(c.DmRpc.RunProxy)
			otaFirmwareM = dmdirect.NewOTAFirmwareManage(c.DmRpc.RunProxy)
			otaTaskM = dmdirect.NewOtaTaskManage(c.DmRpc.RunProxy)
			protocolM = dmdirect.NewProtocolManage(c.DmRpc.RunProxy)
			schemaM = dmdirect.NewSchemaManage(c.DmRpc.RunProxy)
			otaJobM = dmdirect.NewOTAJobManage(c.DmRpc.RunProxy)
			otaModuleM = dmdirect.NewOTAModuleManage(c.DmRpc.RunProxy)
			taskM = dmdirect.NewOTAUpgradeTaskManage(c.DmRpc.RunProxy)
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
			Rule = rule.NewRule(zrpc.MustNewClient(c.UdRpc.Conf))
			Ops = ops.NewOps(zrpc.MustNewClient(c.UdRpc.Conf))
			UserDevice = userdevice.NewUserDevice(zrpc.MustNewClient(c.UdRpc.Conf))
		} else {
			Rule = uddirect.NewRule(c.UdRpc.RunProxy)
			Ops = uddirect.NewOps(c.UdRpc.RunProxy)
			UserDevice = uddirect.NewUserDevice(c.UdRpc.RunProxy)
		}
	}
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
			DataM = datamanage.NewDataManage(zrpc.MustNewClient(c.SysRpc.Conf))
			accessM = accessmanage.NewAccessManage(zrpc.MustNewClient(c.SysRpc.Conf))
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
			DataM = sysdirect.NewData(c.SysRpc.RunProxy)
			accessM = sysdirect.NewAccess(c.SysRpc.RunProxy)
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
			//viddirect.ApiDirectRun()
		}
	}

	if c.VidSip.Enable {
		if c.VidSip.Mode == conf.ClientModeGrpc {
			vidSip = sipmanage.NewSipManage(zrpc.MustNewClient(c.VidSip.Conf))
		} else {
			vidSip = sipdirect.NewSipManage(c.VidSip.RunProxy)
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
		OtaFirmwareM:   otaFirmwareM,
		OtaTaskM:       otaTaskM,
		Ws:             ws.MustNewServer(c.RestConf),
		OtaJobM:        otaJobM,
		OtaModuleM:     otaModuleM,
		TaskM:          taskM,
		SvrClient: SvrClient{
			TenantRpc:      tenantM,
			AppRpc:         appRpc,
			UserRpc:        ur,
			RoleRpc:        ro,
			ModuleRpc:      me,
			LogRpc:         lo,
			AccessRpc:      accessM,
			Timedscheduler: timedSchedule,
			TimedJob:       timedJob,
			VidmgrM:        vidmgrM,
			VidmgrC:        vidmgrC,
			VidmgrS:        vidmgrS,
			SipRpc:         vidSip,

			ProtocolM: protocolM,
			ProjectM:  projectM,
			SchemaM:   schemaM,
			AreaM:     areaM,
			ProductM:  productM,
			DeviceM:   deviceM,
			DeviceA:   deviceA,
			DeviceG:   deviceG,
			DataM:     DataM,

			DeviceMsg:      deviceMsg,
			DeviceInteract: deviceInteract,
			RemoteConfig:   remoteConfig,
			Common:         sysCommon,
			Rule:           Rule,
			UserDevice:     UserDevice,
			Ops:            Ops,
		},
		//OSS:        ossClient,
	}
}
