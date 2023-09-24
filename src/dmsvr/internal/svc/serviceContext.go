package svc

import (
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"os"

	"github.com/i-Things/things/shared/stores"

	"github.com/i-Things/things/shared/caches"

	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/eventBus"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsgManage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/publish/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/deviceDataRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/sdkLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config

	ProjectID *utils.SnowFlake
	AreaID    *utils.SnowFlake
	ProductID *utils.SnowFlake
	DeviceID  *utils.SnowFlake
	GroupID   *utils.SnowFlake
	OssClient *oss.Client

	SchemaRepo     schema.Repo
	SchemaManaRepo deviceMsgManage.SchemaDataRepo
	HubLogRepo     deviceMsgManage.HubLogRepo
	SDKLogRepo     deviceMsgManage.SDKLogRepo
	DataUpdate     dataUpdate.DataUpdate
	Bus            eventBus.Bus
}

func NewServiceContext(c config.Config) *ServiceContext {
	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	ProjectID := utils.NewSnowFlake(nodeID)
	AreaID := utils.NewSnowFlake(nodeID)
	DeviceID := utils.NewSnowFlake(nodeID)
	ProductID := utils.NewSnowFlake(nodeID)
	GroupID := utils.NewSnowFlake(nodeID)

	ccSchemaR := cache.NewSchemaRepo()
	deviceDataR := deviceDataRepo.NewDeviceDataRepo(c.TDengine.DataSource, ccSchemaR.GetSchemaModel)
	hubLogR := hubLogRepo.NewHubLogRepo(c.TDengine.DataSource)
	sdkLogR := sdkLogRepo.NewSDKLogRepo(c.TDengine.DataSource)
	duR, err := dataUpdate.NewDataUpdate(c.Event)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	ossClient := oss.NewOssClient(c.OssConf)
	if ossClient == nil {
		logx.Error("NewOss err")
		os.Exit(-1)
	}
	bus := eventBus.NewEventBus()
	stores.InitConn(c.Database)
	err = relationDB.Migrate()
	if err != nil {
		logx.Error("dmsvr 初始化数据库错误 err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Bus:       bus,
		Config:    c,
		OssClient: ossClient,

		ProjectID: ProjectID,
		AreaID:    AreaID,
		ProductID: ProductID,
		DeviceID:  DeviceID,
		GroupID:   GroupID,

		SchemaRepo:     ccSchemaR,
		SchemaManaRepo: deviceDataR,
		HubLogRepo:     hubLogR,
		SDKLogRepo:     sdkLogR,
		DataUpdate:     duR,
	}
}
