package svc

import (
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/i-Things/share/conf"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgSdkLog"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/service/dmsvr/internal/repo/event/publish/pubApp"
	"github.com/i-Things/things/service/dmsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/service/dmsvr/internal/repo/tdengine/schemaDataRepo"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"

	"gitee.com/i-Things/share/stores"

	"gitee.com/i-Things/share/caches"

	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/config"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"
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
	SchemaRepo     schema.Repo
	SchemaManaRepo msgThing.SchemaDataRepo
	HubLogRepo     msgHubLog.HubLogRepo
	SDKLogRepo     msgSdkLog.SDKLogRepo
	Cache          kv.Store
	ServerMsg      *eventBus.ServerMsg
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

	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	serverMsg, err := eventBus.NewServerMsg(c.Event, c.Name)
	logx.Must(err)
	stores.InitConn(c.Database)
	//err = relationDB.Migrate(c.Database)
	//if err != nil {
	//	logx.Error("dmsvr 初始化数据库错误 err", err)
	//	os.Exit(-1)
	//}
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
		ServerMsg:      serverMsg,
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
	}
}
