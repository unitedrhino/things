package svc

import (
	"gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/common"
	"gitee.com/unitedrhino/core/service/syssvr/client/notifymanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/ops"
	"gitee.com/unitedrhino/core/service/syssvr/client/projectmanage"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicegroup"
	"gitee.com/unitedrhino/things/service/dmsvr/client/deviceinteract"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemsg"
	"gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/userdevice"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/service/dmsvr/dmdirect"
	"gitee.com/unitedrhino/things/service/udsvr/internal/config"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type SvrClient struct {
	ProductM           productmanage.ProductManage
	DeviceInteract     deviceinteract.DeviceInteract
	DeviceMsg          devicemsg.DeviceMsg
	DeviceM            devicemanage.DeviceManage
	DeviceG            devicegroup.DeviceGroup
	TimedM             timedmanage.TimedManage
	NotifyM            notifymanage.NotifyManage
	AreaM              areamanage.AreaManage
	SysCommon          common.Common
	ProjectM           projectmanage.ProjectManage
	DeviceCache        dmExport.DeviceCacheT
	UserShareCache     dmExport.UserShareCacheT
	ProductCache       dmExport.ProductCacheT
	ProductSchemaCache dmExport.SchemaCacheT
	ProjectCache       sysExport.ProjectCacheT
	Ops                ops.Ops
}

type ServiceContext struct {
	Config    config.Config
	FastEvent *eventBus.FastEvent
	Store     kv.Store
	OssClient *oss.Client
	NodeID    int64
	SvrClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedM         timedmanage.TimedManage
		notifyM        notifymanage.NotifyManage
		areaM          areamanage.AreaManage
		projectM       projectmanage.ProjectManage
		deviceM        devicemanage.DeviceManage
		deviceG        devicegroup.DeviceGroup
		productM       productmanage.ProductManage
		userDevice     userdevice.UserDevice
		deviceInteract deviceinteract.DeviceInteract
		deviceMsg      devicemsg.DeviceMsg
		Ops            ops.Ops
	)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	stores.InitConn(c.Database)
	logx.Must(relationDB.Migrate(c.Database))
	timedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
	areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
	Ops = ops.NewOps(zrpc.MustNewClient(c.SysRpc.Conf))
	projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
	notifyM = notifymanage.NewNotifyManage(zrpc.MustNewClient(c.SysRpc.Conf))
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))
		userDevice = userdevice.NewUserDevice(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
		deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
		deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
		deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
		deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
		userDevice = dmdirect.NewUserDevice(c.DmRpc.RunProxy)
	}
	fastEvent, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	dic, err := dmExport.NewDeviceInfoCache(deviceM, fastEvent)
	logx.Must(err)
	pic, err := dmExport.NewProductInfoCache(productM, fastEvent)
	logx.Must(err)
	psc, err := dmExport.NewSchemaInfoCache(productM, fastEvent)
	logx.Must(err)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	udc, err := dmExport.NewUserShareCache(userDevice, fastEvent)
	logx.Must(err)
	projectC, err := sysExport.NewProjectInfoCache(projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf)), fastEvent)
	logx.Must(err)
	return &ServiceContext{
		Config:    c,
		FastEvent: fastEvent,
		Store:     kv.NewStore(c.CacheRedis),
		OssClient: ossClient,
		NodeID:    nodeID,
		SvrClient: SvrClient{
			TimedM:             timedM,
			AreaM:              areaM,
			SysCommon:          common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf)),
			NotifyM:            notifyM,
			ProjectM:           projectM,
			ProjectCache:       projectC,
			ProductM:           productM,
			Ops:                Ops,
			DeviceInteract:     deviceInteract,
			DeviceMsg:          deviceMsg,
			DeviceM:            deviceM,
			DeviceG:            deviceG,
			DeviceCache:        dic,
			ProductCache:       pic,
			UserShareCache:     udc,
			ProductSchemaCache: psc,
		},
	}
}
