package svc

import (
	"gitee.com/unitedrhino/core/service/apisvr/exportMiddleware"
	"gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/log"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	user "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
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
	"gitee.com/unitedrhino/things/service/udsvr/client/rule"
	"gitee.com/unitedrhino/things/service/udsvr/uddirect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
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
	UserM      user.UserManage
	UserC      sysExport.UserCacheT
	AreaC      sysExport.AreaCacheT
	AreaM      areamanage.AreaManage
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

	ur = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
	ro = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
	lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
	areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
	tm = tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
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
	uc, err := sysExport.NewUserInfoCache(ur, serverMsg)
	logx.Must(err)
	areaC, err := sysExport.NewAreaInfoCache(areaM, serverMsg)
	logx.Must(err)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: exportMiddleware.NewCheckTokenWareMiddleware(ur, ro, tm, lo).Handle,
		InitCtxsWare:   ctxs.InitMiddleware,
		OssClient:      ossClient,
		OtaM:           otaM,
		Ws:             ws.MustNewServer(c.RestConf),
		ProductCache:   pc,
		DeviceCache:    dc,
		SvrClient: SvrClient{
			UserM:     ur,
			ProtocolM: protocolM,
			AreaC:     areaC,
			SchemaM:   schemaM,
			ProductM:  productM,
			DeviceM:   deviceM,
			DeviceA:   deviceA,
			UserC:     uc,
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
