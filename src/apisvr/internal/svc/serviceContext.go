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
	"github.com/i-Things/things/src/dmsvr/client/deviceauth"
	"github.com/i-Things/things/src/dmsvr/client/devicegroup"
	"github.com/i-Things/things/src/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/src/dmsvr/client/devicemanage"
	"github.com/i-Things/things/src/dmsvr/client/devicemsg"
	firmwaremanage "github.com/i-Things/things/src/dmsvr/client/firmwaremanage"
	otataskmanage "github.com/i-Things/things/src/dmsvr/client/otataskmanage"
	"github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/client/remoteconfig"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	alarmcenter "github.com/i-Things/things/src/rulesvr/client/alarmcenter"
	scenelinkage "github.com/i-Things/things/src/rulesvr/client/scenelinkage"
	"github.com/i-Things/things/src/rulesvr/ruledirect"
	api "github.com/i-Things/things/src/syssvr/client/api"
	"github.com/i-Things/things/src/syssvr/client/areamanage"
	common "github.com/i-Things/things/src/syssvr/client/common"
	log "github.com/i-Things/things/src/syssvr/client/log"
	menu "github.com/i-Things/things/src/syssvr/client/menu"
	"github.com/i-Things/things/src/syssvr/client/projectmanage"
	role "github.com/i-Things/things/src/syssvr/client/role"
	user "github.com/i-Things/things/src/syssvr/client/user"
	"github.com/i-Things/things/src/syssvr/sysdirect"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
	"github.com/i-Things/things/src/timed/timedjobsvr/timedjobdirect"
	"github.com/i-Things/things/src/timed/timedschedulersvr/client/timedscheduler"
	"github.com/i-Things/things/src/timed/timedschedulersvr/timedschedulerdirect"
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
	UserRpc user.User
	RoleRpc role.Role
	MenuRpc menu.Menu
	LogRpc  log.Log
	ApiRpc  api.Api
	VidmgrM vidmgrinfomanage.VidmgrInfoManage
	VidmgrC vidmgrconfigmanage.VidmgrConfigManage
	VidmgrS vidmgrstreammanage.VidmgrStreamManage

	ProjectM projectmanage.ProjectManage
	AreaM    areamanage.AreaManage
	ProductM productmanage.ProductManage
	DeviceM  devicemanage.DeviceManage
	DeviceA  deviceauth.DeviceAuth
	DeviceG  devicegroup.DeviceGroup

	DeviceMsg      devicemsg.DeviceMsg
	DeviceInteract deviceinteract.DeviceInteract

	RemoteConfig remoteconfig.RemoteConfig
	Common       common.Common

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
	Captcha        *verify.Captcha
	OssClient      *oss.Client
	FirmwareM      firmwaremanage.FirmwareManage
	OtaTaskM       otataskmanage.OtaTaskManage
	FileChan       chan int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		vidmgrM vidmgrinfomanage.VidmgrInfoManage
		vidmgrC vidmgrconfigmanage.VidmgrConfigManage
		vidmgrS vidmgrstreammanage.VidmgrStreamManage

		projectM projectmanage.ProjectManage
		areaM    areamanage.AreaManage
		productM productmanage.ProductManage
		deviceM  devicemanage.DeviceManage
		deviceA  deviceauth.DeviceAuth
		deviceG  devicegroup.DeviceGroup

		deviceMsg      devicemsg.DeviceMsg
		deviceInteract deviceinteract.DeviceInteract
		remoteConfig   remoteconfig.RemoteConfig
		sysCommon      common.Common
		scene          scenelinkage.SceneLinkage
		alarm          alarmcenter.AlarmCenter
		firmwareM      firmwaremanage.FirmwareManage
		otaTaskM       otataskmanage.OtaTaskManage
		timedSchedule  timedscheduler.Timedscheduler
		timedJob       timedmanage.TimedManage
	)
	var ur user.User
	var ro role.Role
	var me menu.Menu
	var lo log.Log
	var ap api.Api

	caches.InitStore(c.CacheRedis)

	ws.StartWsDp(false)
	//var me menu.Menu
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc { //服务模式
			deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DmRpc.Conf))
			productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceA = deviceauth.NewDeviceAuth(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))
			remoteConfig = remoteconfig.NewRemoteConfig(zrpc.MustNewClient(c.DmRpc.Conf))
			firmwareM = firmwaremanage.NewFirmwareManage(zrpc.MustNewClient(c.DmRpc.Conf))
			otaTaskM = otataskmanage.NewOtaTaskManage(zrpc.MustNewClient(c.DmRpc.Conf))
		} else { //直连模式
			deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
			deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
			deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
			productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
			deviceA = dmdirect.NewDeviceAuth(c.DmRpc.RunProxy)
			deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
			remoteConfig = dmdirect.NewRemoteConfig(c.DmRpc.RunProxy)
			firmwareM = dmdirect.NewFirmwareManage(c.DmRpc.RunProxy)
			otaTaskM = dmdirect.NewOtaTaskManage(c.DmRpc.RunProxy)
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
			projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
			areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
			ur = user.NewUser(zrpc.MustNewClient(c.SysRpc.Conf))
			ro = role.NewRole(zrpc.MustNewClient(c.SysRpc.Conf))
			me = menu.NewMenu(zrpc.MustNewClient(c.SysRpc.Conf))
			lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
			ap = api.NewApi(zrpc.MustNewClient(c.SysRpc.Conf))
			sysCommon = common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf))
		} else {
			projectM = sysdirect.NewProjectManage(c.SysRpc.RunProxy)
			areaM = sysdirect.NewAreaManage(c.SysRpc.RunProxy)
			ur = sysdirect.NewUser(c.SysRpc.RunProxy)
			ro = sysdirect.NewRole(c.SysRpc.RunProxy)
			me = sysdirect.NewMenu(c.SysRpc.RunProxy)
			lo = sysdirect.NewLog(c.SysRpc.RunProxy)
			ap = sysdirect.NewApi(c.SysRpc.RunProxy)
			sysCommon = sysdirect.NewCommon(c.SysRpc.RunProxy)
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
		DataAuthWare:   middleware.NewDataAuthWareMiddleware(c).Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware(c, lo).Handle,
		Captcha:        captcha,
		OssClient:      ossClient,
		FirmwareM:      firmwareM,
		OtaTaskM:       otaTaskM,
		SvrClient: SvrClient{
			UserRpc:        ur,
			RoleRpc:        ro,
			MenuRpc:        me,
			LogRpc:         lo,
			ApiRpc:         ap,
			Timedscheduler: timedSchedule,
			TimedJob:       timedJob,
			VidmgrM:        vidmgrM,
			VidmgrC:        vidmgrC,
			VidmgrS:        vidmgrS,

			ProjectM: projectM,
			AreaM:    areaM,
			ProductM: productM,
			DeviceM:  deviceM,
			DeviceA:  deviceA,
			DeviceG:  deviceG,

			DeviceMsg:      deviceMsg,
			DeviceInteract: deviceInteract,
			RemoteConfig:   remoteConfig,
			Common:         sysCommon,
			Scene:          scene,
			Alarm:          alarm,
		},
		//OSS:        ossClient,
	}
}
