package svc

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/coreExport"
	"gitee.com/i-Things/core/service/syssvr/client/areamanage"
	"gitee.com/i-Things/core/service/syssvr/client/datamanage"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	"gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	"gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/domain/tenant"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/service/dmsvr/internal/repo/event/publish/pubApp"
	"github.com/i-Things/things/service/dmsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/repo/tdengine/schemaDataRepo"
	"github.com/i-Things/things/service/dmsvr/internal/repo/tdengine/sendLogRepo"
	"github.com/i-Things/things/service/dmsvr/internal/repo/tdengine/statusLogRepo"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"

	"gitee.com/i-Things/share/stores"

	"gitee.com/i-Things/share/caches"

	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/config"
	"github.com/i-Things/things/service/dmsvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/service/dmsvr/internal/repo/tdengine/sdkLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config

	PubDev pubDev.PubDev
	PubApp pubApp.PubApp

	ProjectID      *utils.SnowFlake
	AreaID         *utils.SnowFlake
	ProductID      *utils.SnowFlake
	DeviceID       *utils.SnowFlake
	GroupID        *utils.SnowFlake
	OssClient      *oss.Client
	TimedM         timedmanage.TimedManage
	SchemaRepo     *caches.Cache[schema.Model]
	SchemaManaRepo msgThing.SchemaDataRepo
	HubLogRepo     deviceLog.HubRepo
	StatusRepo     deviceLog.StatusRepo
	SendRepo       deviceLog.SendRepo
	SDKLogRepo     deviceLog.SDKRepo
	Cache          kv.Store
	DeviceStatus   *cache.DeviceStatus
	FastEvent      *eventBus.FastEvent
	AreaM          areamanage.AreaManage
	UserM          usermanage.UserManage
	DataM          datamanage.DataManage
	ProjectM       projectmanage.ProjectManage
	ProductCache   *caches.Cache[dm.ProductInfo]
	DeviceCache    *caches.Cache[dm.DeviceInfo]
	TenantCache    *caches.Cache[tenant.Info]
	WebHook        *sysExport.Webhook
	UserSubscribe  *coreExport.UserSubscribe
	NodeID         int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedM   timedmanage.TimedManage
		areaM    areamanage.AreaManage
		projectM projectmanage.ProjectManage
	)
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("dmsvr 初始化数据库错误 err", err)
		os.Exit(-1)
	}
	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	ProjectID := utils.NewSnowFlake(nodeID)
	AreaID := utils.NewSnowFlake(nodeID)
	DeviceID := utils.NewSnowFlake(nodeID)
	ProductID := utils.NewSnowFlake(nodeID)
	GroupID := utils.NewSnowFlake(nodeID)
	ca := kv.NewStore(c.CacheRedis)

	hubLogR := hubLogRepo.NewHubLogRepo(c.TSDB)
	sdkLogR := sdkLogRepo.NewSDKLogRepo(c.TSDB)
	statusR := statusLogRepo.NewStatusLogRepo(c.TSDB)
	sendR := sendLogRepo.NewSendLogRepo(c.TSDB)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name)
	logx.Must(err)

	ccSchemaR, err := caches.NewCache(caches.CacheConfig[schema.Model]{
		KeyType:   eventBus.ServerCacheKeyDmSchema,
		FastEvent: serverMsg,
		GetData: func(ctx context.Context, key string) (*schema.Model, error) {
			db := relationDB.NewProductSchemaRepo(ctx)
			dbSchemas, err := db.FindByFilter(ctx, relationDB.ProductSchemaFilter{ProductID: key}, nil)
			if err != nil {
				return nil, err
			}
			schemaModel := relationDB.ToSchemaDo(key, dbSchemas)
			schemaModel.ValidateWithFmt()
			return schemaModel, nil
		},
		Fmt: func(ctx context.Context, key string, data *schema.Model) {
			data.ValidateWithFmt()
		},
		ExpireTime: 10 * time.Minute,
	})
	logx.Must(err)
	deviceDataR := schemaDataRepo.NewDeviceDataRepo(c.TSDB, ccSchemaR.GetData, ca)
	pd, err := pubDev.NewPubDev(c.Event)
	if err != nil {
		logx.Error("NewPubDev err", err)
		os.Exit(-1)
	}
	pa, err := pubApp.NewPubApp(c.Event)
	if err != nil {
		logx.Error("NewPubApp err", err)
		os.Exit(-1)
	}
	timedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
	areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
	userM := usermanage.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
	dataM := datamanage.NewDataManage(zrpc.MustNewClient(c.SysRpc.Conf))
	projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
	tenantCache, err := sysExport.NewTenantInfoCache(tenantmanage.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf)), serverMsg)
	logx.Must(err)
	webHook, err := sysExport.NewTenantOpenWebhook(tenantmanage.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf)), serverMsg)
	logx.Must(err)
	//webHook.Publish(ctxs.BindTenantCode(context.Background(), "default"),
	//	tenantOpenWebhook.CodeDmDeviceConn, application.ConnectMsg{Device: devices.Core{
	//		ProductID:  "123",
	//		DeviceName: "123",
	//	}, Timestamp: time.Now().UnixMilli()})
	return &ServiceContext{
		FastEvent:      serverMsg,
		TenantCache:    tenantCache,
		Config:         c,
		OssClient:      ossClient,
		TimedM:         timedM,
		AreaM:          areaM,
		ProjectM:       projectM,
		PubApp:         pa,
		PubDev:         pd,
		ProjectID:      ProjectID,
		AreaID:         AreaID,
		ProductID:      ProductID,
		DeviceID:       DeviceID,
		GroupID:        GroupID,
		Cache:          ca,
		UserM:          userM,
		DataM:          dataM,
		UserSubscribe:  coreExport.NewUserSubscribe(ca, serverMsg),
		SchemaRepo:     ccSchemaR,
		SchemaManaRepo: deviceDataR,
		DeviceStatus:   cache.NewDeviceStatus(ca),
		HubLogRepo:     hubLogR,
		SDKLogRepo:     sdkLogR,
		StatusRepo:     statusR,
		SendRepo:       sendR,
		WebHook:        webHook,
		NodeID:         nodeID,
	}
}

func (s *ServiceContext) WithDeviceTenant(ctx context.Context, dev devices.Core) context.Context {
	di, err := s.DeviceCache.GetData(ctx, dev.ProductID+":"+dev.DeviceName)
	if err != nil {
		return ctx
	}
	return ctxs.BindTenantCode(ctx, di.TenantCode)
}
