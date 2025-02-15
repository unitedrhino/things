package svc

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/common"
	"gitee.com/unitedrhino/core/service/syssvr/client/datamanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/dictmanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/projectmanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	ws "gitee.com/unitedrhino/share/websocket"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/config"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceBind"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocolTrans"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/event/publish/pubApp"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/event/publish/pubDev"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/abnormalLogRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/hubLogRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/schemaDataRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/sdkLogRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/sendLogRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/statusLogRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"gitee.com/unitedrhino/things/share/topics"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

type ServiceContext struct {
	Config config.Config

	PubDev               pubDev.PubDev
	PubApp               pubApp.PubApp
	ScriptTrans          *protocolTrans.ScriptTrans
	OssClient            *oss.Client
	TimedM               timedmanage.TimedManage
	ProductSchemaRepo    *caches.Cache[schema.Model, string]
	DeviceSchemaRepo     *caches.Cache[schema.Model, devices.Core]
	SchemaManaRepo       msgThing.SchemaDataRepo
	HubLogRepo           deviceLog.HubRepo
	StatusRepo           deviceLog.StatusRepo
	SendRepo             deviceLog.SendRepo
	AbnormalRepo         deviceLog.AbnormalRepo
	SDKLogRepo           deviceLog.SDKRepo
	Cache                kv.Store
	DeviceStatus         *cache.DeviceStatus
	FastEvent            *eventBus.FastEvent
	AreaM                areamanage.AreaManage
	UserM                usermanage.UserManage
	DictM                dictmanage.DictManage
	Common               common.Common
	DataM                datamanage.DataManage
	ProjectM             projectmanage.ProjectManage
	ProductCache         *caches.Cache[dm.ProductInfo, string]
	DeviceCache          *caches.Cache[dm.DeviceInfo, devices.Core]
	UserDeviceShare      *caches.Cache[dm.UserDeviceShareInfo, userShared.UserShareKey]
	UserMultiDeviceShare *caches.Cache[dm.UserDeviceShareMultiInfo, string]
	DeviceBindToken      *caches.Cache[deviceBind.TokenInfo, string]
	TenantCache          sysExport.TenantCacheT
	ProjectCache         sysExport.ProjectCacheT
	AreaCache            sysExport.AreaCacheT
	WebHook              *sysExport.Webhook
	Slot                 sysExport.SlotCacheT
	UserSubscribe        *ws.UserSubscribe
	GatewayCanBind       *cache.GatewayCanBind
	NodeID               int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedM   timedmanage.TimedManage
		areaM    areamanage.AreaManage
		projectM projectmanage.ProjectManage
		Common   common.Common
	)
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("dmsvr 初始化数据库错误 err", err)
		os.Exit(-1)
	}
	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	ca := kv.NewStore(c.CacheRedis)

	hubLogR := hubLogRepo.NewHubLogRepo(c.TSDB)
	abnormalR := abnormalLogRepo.NewAbnormalLogRepo(c.TSDB)
	sdkLogR := sdkLogRepo.NewSDKLogRepo(c.TSDB)
	statusR := statusLogRepo.NewStatusLogRepo(c.TSDB)
	sendR := sendLogRepo.NewSendLogRepo(c.TSDB)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)

	getProductSchemaModel, err := caches.NewCache(caches.CacheConfig[schema.Model, string]{
		KeyType:   topics.ServerCacheKeyDmProductSchema,
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
		Fmt: func(ctx context.Context, key string, data *schema.Model) *schema.Model {
			data.ValidateWithFmt()
			return data
		},
		ExpireTime: 10 * time.Minute,
	})
	logx.Must(err)
	getDeviceSchemaModel, err := caches.NewCache(caches.CacheConfig[schema.Model, devices.Core]{
		KeyType:   topics.ServerCacheKeyDmDeviceSchema,
		FastEvent: serverMsg,
		GetData: func(ctx context.Context, key devices.Core) (*schema.Model, error) {
			db := relationDB.NewDeviceSchemaRepo(ctx)
			dbSchemas, err := db.FindByFilter(ctx, relationDB.DeviceSchemaFilter{
				ProductID: key.ProductID, DeviceName: key.DeviceName}, nil)
			if err != nil {
				return nil, err
			}
			schemaModel := relationDB.ToDeviceSchemaDo(key.ProductID, dbSchemas)
			schemaModel.ValidateWithFmt()
			return schemaModel, nil
		},
		Fmt: func(ctx context.Context, key devices.Core, data *schema.Model) *schema.Model {
			ps, err := getProductSchemaModel.GetData(ctx, key.ProductID)
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
			newOne := data.Copy().Aggregation(ps)
			newOne.ValidateWithFmt()
			return newOne
		},
		ExpireTime: 20 * time.Minute,
	})
	getProductSchemaModel.AddNotifySlot(func(ctx context.Context, keyB []byte) {
		var pKey devices.Core
		json.Unmarshal(keyB, &pKey)
		getDeviceSchemaModel.DeleteByFunc(func(key string) bool {
			ck := devices.Core{}
			json.Unmarshal([]byte(key), &ck)
			if ck.ProductID == pKey.ProductID {
				return true
			}
			return false
		})
	})
	logx.Must(err)
	deviceDataR := schemaDataRepo.NewDeviceDataRepo(c.TSDB, func(ctx context.Context, productID devices.Core) (*schema.Model, error) {
		return getProductSchemaModel.GetData(ctx, productID.ProductID)
	}, getDeviceSchemaModel.GetData, ca)
	err = deviceDataR.Init(context.Background())
	logx.Must(err)
	pd, err := pubDev.NewPubDev(serverMsg)
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
	dictM := dictmanage.NewDictManage(zrpc.MustNewClient(c.SysRpc.Conf))
	dataM := datamanage.NewDataManage(zrpc.MustNewClient(c.SysRpc.Conf))
	projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
	Common = common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf))
	tenantCache, err := sysExport.NewTenantInfoCache(tenantmanage.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf)), serverMsg)
	logx.Must(err)
	webHook, err := sysExport.NewTenantOpenWebhook(tenantmanage.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf)), serverMsg)
	logx.Must(err)
	//webHook.Publish(ctxs.BindTenantCode(context.Background(), "default"),
	//	tenantOpenWebhook.CodeDmDeviceConn, application.ConnectMsg{Device: devices.Core{
	//		ProductID:  "123",
	//		DeviceName: "123",
	//	}, Timestamp: time.Now().UnixMilli()})
	Slot, err := sysExport.NewSlotCache(common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf)))
	logx.Must(err)
	projectC, err := sysExport.NewProjectInfoCache(projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf)), serverMsg)
	logx.Must(err)
	areaC, err := sysExport.NewAreaInfoCache(areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf)), serverMsg)
	logx.Must(err)
	return &ServiceContext{
		FastEvent:         serverMsg,
		TenantCache:       tenantCache,
		Config:            c,
		OssClient:         ossClient,
		TimedM:            timedM,
		AreaM:             areaM,
		ProjectM:          projectM,
		PubApp:            pa,
		PubDev:            pd,
		Cache:             ca,
		UserM:             userM,
		DictM:             dictM,
		Common:            Common,
		DataM:             dataM,
		UserSubscribe:     ws.NewUserSubscribe(ca, serverMsg),
		ProductSchemaRepo: getProductSchemaModel,
		DeviceSchemaRepo:  getDeviceSchemaModel,
		SchemaManaRepo:    deviceDataR,
		DeviceStatus:      cache.NewDeviceStatus(ca),
		GatewayCanBind:    cache.NewGatewayCanBind(ca),
		HubLogRepo:        hubLogR,
		SDKLogRepo:        sdkLogR,
		AbnormalRepo:      abnormalR,
		StatusRepo:        statusR,
		SendRepo:          sendR,
		WebHook:           webHook,
		NodeID:            nodeID,
		Slot:              Slot,
		ProjectCache:      projectC,
		ScriptTrans:       protocolTrans.NewScriptTrans(),
		AreaCache:         areaC,
	}
}

func (s *ServiceContext) WithDeviceTenant(ctx context.Context, dev devices.Core) context.Context {
	di, err := s.DeviceCache.GetData(ctx, dev)
	if err != nil {
		return ctx
	}
	return ctxs.BindTenantCode(ctx, di.TenantCode, di.ProjectID)
}

//func test() {
//	i := interp.New(interp.Options{})
//
//	i.Use(stdlib.Symbols)
//
//	_, err := i.Eval(`import "ur/ur"`)
//	if err != nil {
//		panic(err)
//	}
//
//	funv, err := i.Eval(`ur.Hello`)
//	if err != nil {
//		panic(err)
//	}
//	fn, ok := funv.Interface().(func(context.Context, string))
//	if !ok {
//		panic(ok)
//	}
//	fn(context.Background(), "hello world")
//}
