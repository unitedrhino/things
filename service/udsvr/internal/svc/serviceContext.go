package svc

import (
	"gitee.com/i-Things/core/service/syssvr/client/areamanage"
	"gitee.com/i-Things/core/service/syssvr/client/notifymanage"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/client/devicegroup"
	"github.com/i-Things/things/service/dmsvr/client/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/devicemsg"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/dmsvr/dmdirect"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/i-Things/things/service/udsvr/internal/config"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
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
	ProjectM           projectmanage.ProjectManage
	DeviceCache        *caches.Cache[dm.DeviceInfo]
	ProductSchemaCache *caches.Cache[schema.Model]
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
		deviceInteract deviceinteract.DeviceInteract
		deviceMsg      devicemsg.DeviceMsg
	)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	stores.InitConn(c.Database)
	logx.Must(relationDB.Migrate(c.Database))
	timedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
	areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
	projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
	notifyM = notifymanage.NewNotifyManage(zrpc.MustNewClient(c.SysRpc.Conf))
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceMsg = devicemsg.NewDeviceMsg(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DmRpc.Conf))
		deviceG = devicegroup.NewDeviceGroup(zrpc.MustNewClient(c.DmRpc.Conf))

	} else {
		productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
		deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
		deviceMsg = dmdirect.NewDeviceMsg(c.DmRpc.RunProxy)
		deviceInteract = dmdirect.NewDeviceInteract(c.DmRpc.RunProxy)
		deviceG = dmdirect.NewDeviceGroup(c.DmRpc.RunProxy)
	}
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	dic, err := dmExport.NewDeviceInfoCache(deviceM, serverMsg)
	logx.Must(err)
	psc, err := dmExport.NewSchemaInfoCache(productM, serverMsg)
	logx.Must(err)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:    c,
		FastEvent: serverMsg,
		Store:     kv.NewStore(c.CacheRedis),
		OssClient: ossClient,
		NodeID:    nodeID,
		SvrClient: SvrClient{
			TimedM:             timedM,
			AreaM:              areaM,
			NotifyM:            notifyM,
			ProjectM:           projectM,
			ProductM:           productM,
			DeviceInteract:     deviceInteract,
			DeviceMsg:          deviceMsg,
			DeviceM:            deviceM,
			DeviceG:            deviceG,
			DeviceCache:        dic,
			ProductSchemaCache: psc,
		},
	}
}
