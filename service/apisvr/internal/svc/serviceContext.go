package svc

import (
	"os"

	"gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/dictmanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/log"
	operLog "gitee.com/unitedrhino/core/service/syssvr/client/log"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	user "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/core/share/middlewares"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/utils"
	ws "gitee.com/unitedrhino/share/websocket"
	"gitee.com/unitedrhino/things/service/apisvr/internal/config"
	"gitee.com/unitedrhino/things/service/dgsvr/client/deviceauth"
	"gitee.com/unitedrhino/things/service/dgsvr/dgdirect"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicegroup"
	"gitee.com/unitedrhino/things/service/dmsvr/client/deviceinteract"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemsg"
	"gitee.com/unitedrhino/things/service/dmsvr/client/otamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/protocolmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/remoteconfig"
	"gitee.com/unitedrhino/things/service/dmsvr/client/schemamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/userdevice"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/service/dmsvr/dmdirect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
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

	UserDevice userdevice.UserDevice
	UserM      user.UserManage
	UserC      sysExport.UserCacheT
	AreaC      sysExport.AreaCacheT
	AreaM      areamanage.AreaManage
	DictM      dictmanage.DictManage
	RoleRpc    role.RoleManage
	TenantRpc  tenant.TenantManage
	LogRpc     operLog.Log
}

type ServiceContext struct {
	SvrClient
	Ws             *ws.Server
	Config         config.Config
	InitCtxsWare   rest.Middleware
	CheckTokenWare rest.Middleware
	OssClient      *oss.Client
	OtaM           otamanage.OtaManage
	ProductCache   dmExport.ProductCacheT
	DeviceCache    dmExport.DeviceCacheT
	UserShareCache dmExport.UserShareCacheT
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
		areaM          areamanage.AreaManage
	)
	var ur user.UserManage
	var ro role.RoleManage
	var tm tenant.TenantManage
	var lo log.Log

	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	//var me menu.Menu
	if c.DmRpc.Enable {
		if c.DmRpc.Mode == conf.ClientModeGrpc { //服务模式
			deviceMsg = devicemsg.NewDeviceMsg(c.DmRpc.MustNewClient())
			deviceInteract = deviceinteract.NewDeviceInteract(c.DmRpc.MustNewClient())
			productM = productmanage.NewProductManage(c.DmRpc.MustNewClient())
			deviceM = devicemanage.NewDeviceManage(c.DmRpc.MustNewClient())
			deviceG = devicegroup.NewDeviceGroup(c.DmRpc.MustNewClient())
			remoteConfig = remoteconfig.NewRemoteConfig(c.DmRpc.MustNewClient())
			otaM = otamanage.NewOtaManage(c.DmRpc.MustNewClient())
			protocolM = protocolmanage.NewProtocolManage(c.DmRpc.MustNewClient())
			schemaM = schemamanage.NewSchemaManage(c.DmRpc.MustNewClient())
			UserDevice = userdevice.NewUserDevice(c.DmRpc.MustNewClient())

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
	udc, err := dmExport.NewUserShareCache(UserDevice, serverMsg)
	logx.Must(err)
	if c.DgRpc.Enable {
		if c.DgRpc.Mode == conf.ClientModeGrpc { //服务模式
			deviceA = deviceauth.NewDeviceAuth(c.DgRpc.MustNewClient())

		} else { //直连模式
			deviceA = dgdirect.NewDeviceAuth(c.DgRpc.RunProxy)
		}
	}
	//if c.UdRpc.Enable {
	//	if c.UdRpc.Mode == conf.ClientModeGrpc {
	//		Rule = rule.NewRule(c.UdRpc.MustNewClient())
	//	} else {
	//		Rule = uddirect.NewRule(c.UdRpc.RunProxy)
	//	}
	//}
	dictM := dictmanage.NewDictManage(c.SysRpc.MustNewClient())
	ur = user.NewUserManage(c.SysRpc.MustNewClient())
	ro = role.NewRoleManage(c.SysRpc.MustNewClient())
	lo = log.NewLog(c.SysRpc.MustNewClient())
	areaM = areamanage.NewAreaManage(c.SysRpc.MustNewClient())
	tm = tenant.NewTenantManage(c.SysRpc.MustNewClient())
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}

	pc, err := dmExport.NewProductInfoCache(productM, serverMsg)
	logx.Must(err)
	dc, err := dmExport.NewDeviceInfoCache(deviceM, serverMsg)
	logx.Must(err)
	uc, err := sysExport.NewUserInfoCache(ur, serverMsg)
	logx.Must(err)
	areaC, err := sysExport.NewAreaInfoCache(areaM, serverMsg)
	logx.Must(err)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: middlewares.NewCheckTokenWareMiddleware(ur, ro, tm, lo).Handle,
		InitCtxsWare:   middlewares.InitMiddleware,
		OssClient:      ossClient,
		OtaM:           otaM,
		Ws:             ws.MustNewServer(c.RestConf),
		ProductCache:   pc,
		DeviceCache:    dc,
		UserShareCache: udc,
		SvrClient: SvrClient{
			UserM:          ur,
			ProtocolM:      protocolM,
			AreaC:          areaC,
			SchemaM:        schemaM,
			ProductM:       productM,
			DeviceM:        deviceM,
			DeviceA:        deviceA,
			UserC:          uc,
			DeviceG:        deviceG,
			AreaM:          areaM,
			DictM:          dictM,
			DeviceMsg:      deviceMsg,
			DeviceInteract: deviceInteract,
			RemoteConfig:   remoteConfig,
			UserDevice:     UserDevice,
			RoleRpc:        ro,
			TenantRpc:      tm,
			LogRpc:         lo,
		},
		//OSS:        ossClient,
	}
}
