package svc

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgSdkLog"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/publish/pubApp"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/schemaDataRepo"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
	"github.com/i-Things/things/src/timed/timedjobsvr/timedjobdirect"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"

	"github.com/i-Things/things/shared/stores"

	"github.com/i-Things/things/shared/caches"

	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/eventBus"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/publish/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/sdkLogRepo"
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
	SchemaRepo     schema.Repo
	SchemaManaRepo msgThing.SchemaDataRepo
	HubLogRepo     msgHubLog.HubLogRepo
	SDKLogRepo     msgSdkLog.SDKLogRepo
	DataUpdate     dataUpdate.DataUpdate
	Cache          kv.Store
	Bus            eventBus.Bus
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedM timedmanage.TimedManage
	)
	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	ProjectID := utils.NewSnowFlake(nodeID)
	AreaID := utils.NewSnowFlake(nodeID)
	DeviceID := utils.NewSnowFlake(nodeID)
	ProductID := utils.NewSnowFlake(nodeID)
	GroupID := utils.NewSnowFlake(nodeID)
	ca := kv.NewStore(c.CacheRedis)
	ccSchemaR := cache.NewSchemaRepo()
	deviceDataR := schemaDataRepo.NewDeviceDataRepo(c.TSDB, ccSchemaR.GetSchemaModel, ca)
	hubLogR := hubLogRepo.NewHubLogRepo(c.TSDB)
	sdkLogR := sdkLogRepo.NewSDKLogRepo(c.TSDB)
	duR, err := dataUpdate.NewDataUpdate(c.Event)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}

	bus := eventBus.NewEventBus()
	stores.InitConn(c.Database)
	err = relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("dmsvr 初始化数据库错误 err", err)
		os.Exit(-1)
	}
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
	if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
		timedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
	} else {
		timedM = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
	}
	return &ServiceContext{
		Bus:            bus,
		Config:         c,
		OssClient:      ossClient,
		TimedM:         timedM,
		PubApp:         pa,
		PubDev:         pd,
		ProjectID:      ProjectID,
		AreaID:         AreaID,
		ProductID:      ProductID,
		DeviceID:       DeviceID,
		GroupID:        GroupID,
		Cache:          ca,
		SchemaRepo:     ccSchemaR,
		SchemaManaRepo: deviceDataR,
		HubLogRepo:     hubLogR,
		SDKLogRepo:     sdkLogR,
		DataUpdate:     duR,
	}
}
