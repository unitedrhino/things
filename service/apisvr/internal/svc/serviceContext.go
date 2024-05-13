package svc

import (
	"gitee.com/i-Things/core/service/apisvr/export"
	"gitee.com/i-Things/core/service/syssvr/client/areamanage"
	"gitee.com/i-Things/core/service/syssvr/client/log"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/sysdirect"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/share/verify"
	ws "gitee.com/i-Things/share/websocket"
	"github.com/i-Things/things/service/apisvr/internal/config"
	"github.com/i-Things/things/service/apisvr/internal/middleware"
	"github.com/i-Things/things/service/dgsvr/client/deviceauth"
	"github.com/i-Things/things/service/dgsvr/dgdirect"
	"github.com/i-Things/things/service/dmsvr/client/devicegroup"
	"github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/devicemsg"
	"github.com/i-Things/things/service/dmsvr/client/otamanage"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/client/protocolmanage"
	"github.com/i-Things/things/service/dmsvr/client/remoteconfig"
	"github.com/i-Things/things/service/dmsvr/client/schemamanage"
	"github.com/i-Things/things/service/dmsvr/client/userdevice"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/dmsvr/dmdirect"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/i-Things/things/service/rulesvr/client/alarmcenter"
	"github.com/i-Things/things/service/rulesvr/client/scenelinkage"
	"github.com/i-Things/things/service/udsvr/client/rule"
	"github.com/i-Things/things/service/udsvr/uddirect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

type SvrClient struct {
	ProtocolM protocolmanage.ProtocolManage
	ProductM  productmanage.ProductManage
	SchemaM   schemamanage.SchemaManage
	DeviceM   devicemanage.DeviceManage
	DeviceA   deviceauth.DeviceAuth
	DeviceG   devicegroup.DeviceGroup

	DeviceMsg      devicemsg.DeviceMsg
	DeviceInteract deviceinteract.DeviceInteract

	RemoteConfig remoteconfig.RemoteConfig

	Rule       rule.Rule
	UserDevice userdevice.UserDevice
	Scene      scenelinkage.SceneLinkage
	Alarm      alarmcenter.AlarmCenter
	UserM      user.UserManage
	ProjectM   projectmanage.ProjectManage
	AreaM      areamanage.AreaManage
}

type ServiceContext struct {
	SvrClient
	Ws             *ws.Server
	Config         config.Config
	InitCtxsWare   rest.Middleware
	SetupWare      rest.Middleware
	CheckTokenWare rest.Middleware
	DataAuthWare   rest.Middleware
	TeardownWare   rest.Middleware
	CheckApiWare   rest.Middleware
	Captcha        *verify.Captcha
	OssClient      *oss.Client
	OtaM           otamanage.OtaManage
	ProductCache   *caches.Cache[dm.ProductInfo]
	DeviceCache    *caches.Cache[dm.DeviceInfo]
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		schemaM schemamanage.SchemaManage

		protocolM protocolmanage.ProtocolManage
		productM  productmanage.ProductManage
		deviceM   devicemanage.DeviceManage
		deviceA   deviceauth.DeviceAuth
		deviceG   devicegroup.DeviceGroup

		deviceMsg      devicemsg.DeviceMsg
		deviceInteract deviceinteract.DeviceInteract
		remoteConfig   remoteconfig.RemoteConfig
		otaM           otamanage.OtaManage
		UserDevice     userdevice.UserDevice
		Rule           rule.Rule
		projectM       projectmanage.ProjectManage
		areaM          areamanage.AreaManage
	)
	var ur user.UserManage
	var ro role.RoleManage
	var tm tenant.TenantManage
	var lo log.Log

	caches.InitStore(c.CacheRedis)

	//var me menu.Menu
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc { //服务模式
			deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DmRpc.Conf))
			productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
			deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))
			remoteConfig = remoteconfig.NewRemoteConfig(zrpc.MustNewClient(c.DmRpc.Conf))
			otaM = otamanage.NewOtaManage(zrpc.MustNewClient(c.DmRpc.Conf))
			protocolM = protocolmanage.NewProtocolManage(zrpc.MustNewClient(c.DmRpc.Conf))
			schemaM = schemamanage.NewSchemaManage(zrpc.MustNewClient(c.DmRpc.Conf))
			UserDevice = userdevice.NewUserDevice(zrpc.MustNewClient(c.UdRpc.Conf))

		} else { //直连模式
			deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
			deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
			deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
			productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
			deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
			remoteConfig = dmdirect.NewRemoteConfig(c.DmRpc.RunProxy)
			otaM = dmdirect.NewOtaManage(c.DmRpc.RunProxy)
			protocolM = dmdirect.NewProtocolManage(c.DmRpc.RunProxy)
			schemaM = dmdirect.NewSchemaManage(c.DmRpc.RunProxy)
			UserDevice = dmdirect.NewUserDevice(c.DmRpc.RunProxy)

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
		} else {
			Rule = uddirect.NewRule(c.UdRpc.RunProxy)
		}
	}
	if c.SysRpc.Enable {
		if c.SysRpc.Mode == conf.ClientModeGrpc {
			ur = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
			ro = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
			lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
			areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
			tm = tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
			projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
		} else {
			ur = sysdirect.NewUser(c.SysRpc.RunProxy)
			ro = sysdirect.NewRole(c.SysRpc.RunProxy)
			lo = sysdirect.NewLog(c.SysRpc.RunProxy)
			areaM = sysdirect.NewAreaManage(c.SysRpc.RunProxy)
			tm = sysdirect.NewTenantManage(c.SysRpc.RunProxy)
			projectM = sysdirect.NewProjectManage(c.SysRpc.RunProxy)
		}
	}

	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	pc, err := dmExport.NewProductInfoCache(productM, serverMsg)
	logx.Must(err)
	dc, err := dmExport.NewDeviceInfoCache(deviceM, serverMsg)
	logx.Must(err)
	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:         c,
		SetupWare:      middleware.NewSetupWareMiddleware(c, lo).Handle,
		CheckTokenWare: export.NewCheckTokenWareMiddleware(ur, ro, tm).Handle,
		InitCtxsWare:   ctxs.InitMiddleware,
		DataAuthWare:   middleware.NewDataAuthWareMiddleware(c).Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware(c, lo).Handle,
		CheckApiWare:   middleware.NewCheckApiWareMiddleware().Handle,
		Captcha:        captcha,
		OssClient:      ossClient,
		OtaM:           otaM,
		Ws:             ws.MustNewServer(c.RestConf),
		ProductCache:   pc,
		DeviceCache:    dc,
		SvrClient: SvrClient{
			UserM:     ur,
			ProjectM:  projectM,
			ProtocolM: protocolM,
			SchemaM:   schemaM,
			ProductM:  productM,
			DeviceM:   deviceM,
			DeviceA:   deviceA,
			DeviceG:   deviceG,
			AreaM:     areaM,

			DeviceMsg:      deviceMsg,
			DeviceInteract: deviceInteract,
			RemoteConfig:   remoteConfig,
			Rule:           Rule,
			UserDevice:     UserDevice,
		},
		//OSS:        ossClient,
	}
}
